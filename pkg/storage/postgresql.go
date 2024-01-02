// Package storage store resource information
package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"

	_ "github.com/lib/pq" // postgresql driver
)

// PostgreSQL store storage configuration for PostgreSQL database
type PostgreSQL struct {
	DB *sql.DB
}

// Prepare to prepare the database
func (p *PostgreSQL) Prepare() error {
	sqlString := `
    CREATE TABLE IF NOT EXISTS resources (
        id SERIAL PRIMARY KEY,
        source text NOT NULL,
        kind TEXT NOT NULL,
        uuid TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL,
        external_id TEXT NOT NULL,
        version INT NOT NULL DEFAULT '1',
        discovered_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
	CREATE UNIQUE INDEX IF NOT EXISTS idx_resources_on_conflict ON resources (kind, uuid, external_id, version);

    CREATE TABLE IF NOT EXISTS metadata (
        resource_id INTEGER REFERENCES resources(id) ON DELETE CASCADE,
        key TEXT NOT NULL,
        value TEXT NOT NULL,
        PRIMARY KEY(resource_id, key, value)
    );

    CREATE TABLE IF NOT EXISTS relations (
        resource_id INTEGER REFERENCES resources(id) ON DELETE CASCADE,
        related_resource_id INTEGER REFERENCES resources(id),
        PRIMARY KEY(resource_id, related_resource_id)
    );

	CREATE INDEX IF NOT EXISTS metadata_key_idx ON metadata (key);
	CREATE INDEX IF NOT EXISTS metadata_value_idx ON metadata (value);
`

	_, err := p.DB.Exec(sqlString)
	if err != nil {
		return fmt.Errorf("failed to create tables in database: %v", err)
	}

	return nil
}

// Persist store resources
func (p *PostgreSQL) Persist(resources []resource.Resource) error {
	return p.runInTransaction(func(tx *sql.Tx) error {
		// Prepare the SQL statements
		resourcesStmt, err := tx.Prepare("INSERT INTO resources (kind, uuid, name, external_id, source, version) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (kind, uuid, external_id, version) DO UPDATE SET kind = excluded.kind, uuid = excluded.uuid, name = excluded.name, source = excluded.source, version = excluded.version RETURNING id")
		if err != nil {
			return fmt.Errorf("failed to prepare resources statement: %w", err)
		}
		defer func(resourcesStmt *sql.Stmt) {
			err := resourcesStmt.Close()
			if err != nil {
				return
			}
		}(resourcesStmt)

		metadataStmt, err := tx.Prepare("INSERT INTO metadata (resource_id, key, value) VALUES ($1, $2, $3) ON CONFLICT (resource_id, key, value) DO UPDATE SET value = excluded.value")
		if err != nil {
			return fmt.Errorf("failed to prepare metadata statement: %w", err)
		}
		defer func(metadataStmt *sql.Stmt) {
			err := metadataStmt.Close()
			if err != nil {
				return
			}
		}(metadataStmt)

		// Loop through the resources and insert or update them into the database
		for _, res := range resources {
			// Insert or update the resource
			var id int
			err := resourcesStmt.QueryRow(res.GetKind(), res.GetUUID(), res.GetName(), res.GetExternalID(), res.GetScanner(), res.GetVersion()).Scan(&id)
			if err != nil {
				return fmt.Errorf("failed to insert or update resource: %w", err)
			}

			// Insert or update the metadata
			data := res.GetMetaData()
			for _, meta := range data.Get() {
				_, err = metadataStmt.Exec(id, meta.Key, meta.Value)
				if err != nil {
					return fmt.Errorf("failed to insert or update metadata: %w", err)
				}
			}
		}

		return nil
	})
}

// Find resources based on criteria
func (p *PostgreSQL) Find(filter ResourceFilter) ([]resource.Resource, error) {
	var resources []resource.Resource

	// Prepare the query
	q := &Query{}
	if filter.Kind != "" {
		q.AddFilter("r.kind", "=", filter.Kind)
	}
	if filter.UUID != "" {
		q.AddFilter("r.uuid", "=", filter.UUID)
	}
	if filter.Name != "" {
		q.AddFilter("r.name", "=", filter.Name)
	}
	if filter.ExternalID != "" {
		q.AddFilter("r.external_id", "=", filter.ExternalID)
	}

	// Build the query
	query, args := q.Build()

	// Execute the query
	rows, err := p.DB.Query(query, args...)
	if err != nil {
		return resources, fmt.Errorf("failed to execute query to find resource: %w", err)
	}
	defer rows.Close()

	// Parse the results
	for rows.Next() {
		var kind, uuid, name, externalID, source string
		var version int
		var discoveredAt time.Time
		var metaJson []byte // For the aggregated metadata JSON

		err := rows.Scan(
			&kind,
			&uuid,
			&name,
			&externalID,
			&source,
			&version,
			&discoveredAt,
			&metaJson, // Scan the JSON data
		)
		if err != nil {
			return resources, fmt.Errorf("failed to scan row to fetch the resource result: %w", err)
		}

		// Create a resource object and use setters
		res := resource.NewResource(kind, name, externalID, source, version)
		res.SetUUID(uuid)
		res.SetFetchedAt(discoveredAt)

		var metaData map[string]string
		if err := json.Unmarshal(metaJson, &metaData); err != nil {
			return resources, fmt.Errorf("failed to unmarshal metadata JSON: %w", err)
		}

		// Assuming there's a method to set all metadata at once
		res.AddMetaData(metaData)

		resources = append(resources, res)
	}

	return resources, nil
}

// StoreRelations will go through all resources in the database and based on criteria it will insert relationship to "relations" table
// StoreRelations will go through all resources in the database and based on criteria it will insert relationship to "relations" table
func (p *PostgreSQL) StoreRelations(relation config.Relation) error {
	return p.runInTransaction(func(tx *sql.Tx) error {
		// Prepare the SQL statement to insert relationships into the relations table
		relationsStmt, err := tx.Prepare("INSERT INTO relations (resource_id, related_resource_id) VALUES ($1, $2)")
		if err != nil {
			return fmt.Errorf("failed to prepare relations statement: %w", err)
		}
		defer func(relationsStmt *sql.Stmt) {
			err := relationsStmt.Close()
			if err != nil {
				return
			}
		}(relationsStmt)

		for _, rc := range relation.RelationCriteria {
			err = p.storeRelationMatrix(rc, relationsStmt)
			if err != nil {
				return fmt.Errorf("failed to store relation matrix: %w", err)
			}
		}

		return nil
	})
}

func (p *PostgreSQL) storeRelationMatrix(rc config.RelationCriteria, relationsStmt *sql.Stmt) error {
	relationMatrix, err := p.analyzeRelationMatrix(rc)
	if err != nil {
		return fmt.Errorf("failed to analyze relation matrix: %w", err)
	}

	// Insert relationships for the matching resources
	for _, c := range relationMatrix {
		for k, v := range c {
			if _, err := relationsStmt.Exec(k, v); err != nil {
				return fmt.Errorf("failed to insert relation: %w", err)
			}
		}
	}
	return nil
}

func (p *PostgreSQL) analyzeRelationMatrix(relCriteria config.RelationCriteria) ([]map[string]string, error) {
	resourceStmt := `SELECT STRING_AGG(r.id::text, ',') AS resource_ids
FROM resources r
         LEFT JOIN metadata m ON r.id = m.resource_id
WHERE m.key = $1 AND m.value = $2 AND r.kind = $3;`

	var relateToIds string
	err := p.DB.QueryRow(resourceStmt, relCriteria.Target.MetaKey, relCriteria.Target.MetaValue, relCriteria.Target.Kind).Scan(&relateToIds)
	if err != nil {
		return nil, err
	}

	var resourceForBuildRelations string
	err = p.DB.QueryRow(resourceStmt, relCriteria.Source.MetaKey, relCriteria.Source.MetaValue, relCriteria.Source.Kind).Scan(&resourceForBuildRelations)
	if err != nil {
		return nil, err
	}

	relationMatrix := p.generateRelationMatrix(
		strings.Split(relateToIds, ","),
		strings.Split(resourceForBuildRelations, ","),
	)
	return relationMatrix, nil
}

func (p *PostgreSQL) generateRelationMatrix(relatedToIds []string, resourceForBuildRelations []string) []map[string]string {
	var matrix []map[string]string
	for _, val1 := range relatedToIds {
		for _, val2 := range resourceForBuildRelations {
			m := map[string]string{
				val2: val1,
			}
			matrix = append(matrix, m)
		}
	}
	return matrix
}

// GetResources fetch all resources from storage
func (p *PostgreSQL) GetResources() ([]resource.Resource, error) {
	var resources []resource.Resource

	// Query all resources
	rows, err := p.DB.Query("SELECT kind, uuid, name, external_id FROM resources")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query to get resources: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	// Loop through the result set and create Resource objects
	for rows.Next() {
		var k string
		var u string
		var n string
		var eid string
		if err := rows.Scan(&k, &u, &n, &eid); err != nil {
			return nil, err
		}
		r := resource.NewResource(k, n, eid, "", 1)
		r.SetUUID(u)
		resources = append(resources, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan row to fetch the resource result: %w", err)
	}

	return resources, nil
}

// GetRelations fetch all the relations between resources
func (p *PostgreSQL) GetRelations() ([]map[string]string, error) {
	var relations []map[string]string

	// Query all relations
	rows, err := p.DB.Query("SELECT r.uuid as resource_uuid, r2.uuid as related_resource_uuid FROM relations left join resources r on r.id = relations.resource_id left join resources r2 on r2.id = relations.related_resource_id")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query to get relations: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	// Loop through the result set and create maps with resource_id and related_resource_id
	for rows.Next() {
		r := map[string]string{}
		var resourceUUID, relatedResourceUUID string
		if err := rows.Scan(&resourceUUID, &relatedResourceUUID); err != nil {
			return nil, fmt.Errorf("failed to scan row to fetch the relation result: %w", err)
		}
		r[resourceUUID] = relatedResourceUUID
		relations = append(relations, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan row to fetch the relation result: %w", err)
	}

	return relations, nil
}

func (p *PostgreSQL) runInTransaction(f func(tx *sql.Tx) error) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			return
		}
		err = tx.Commit()
	}()

	return f(tx)
}

// GetNextVersionForResource get next version for a resource
func (p *PostgreSQL) GetNextVersionForResource(source, kind string) (int, error) {
	var version int
	err := p.DB.QueryRow("SELECT COALESCE(max(version), 0) + 1 FROM resources WHERE source = $1 AND kind = $2", source, kind).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("failed to get next version for resource: %w", err)
	}

	return version, nil
}

// CleanupOldVersion cleanup old version of resources
func (p *PostgreSQL) CleanupOldVersion(source, kind string) (int64, error) {
	// Construct the SQL statement
	query := `
	DELETE FROM resources
	WHERE source = $1
	  AND kind = $2
	  AND version NOT IN (
		SELECT DISTINCT version
		FROM resources
		WHERE source = $1
		  AND kind = $2
		ORDER BY version DESC
		LIMIT $3
	  );
	`

	// Execute the query
	result, err := p.DB.Exec(query, source, kind, noOfVersionToKeep)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup the resources: %v", err)
	}

	// Get the count of affected rows
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve affected rows count: %v", err)
	}

	return affectedRows, nil
}

// GetResourceCount get resource count
func (p *PostgreSQL) GetResourceCount() ([]ResourceCount, error) {
	const query = `
	SELECT r.source, r.kind, COUNT(r.id) as total_count
	FROM resources r
	JOIN (
	    SELECT source, kind, MAX(version) as latest_version
	    FROM resources
	    GROUP BY kind, source
	) sub ON r.kind = sub.kind AND r.version = sub.latest_version AND r.source = sub.source
	GROUP BY r.kind, r.source
	ORDER BY r.kind;
	`

	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []ResourceCount
	for rows.Next() {
		var rc ResourceCount
		err := rows.Scan(&rc.Source, &rc.Kind, &rc.TotalCount)
		if err != nil {
			return nil, err
		}
		results = append(results, rc)
	}

	return results, nil
}

// GetResourceCountByMetaData get filtered resource count
func (p *PostgreSQL) GetResourceCountByMetaData() ([]MetadataCount, error) {
	const query = `
	SELECT r.source, r.kind, m.key, m.value, COUNT(distinct m.resource_id) as total_count
	FROM metadata m
	JOIN resources r ON m.resource_id = r.id
	JOIN (
	    SELECT source, kind, MAX(version) as latest_version
	    FROM resources
	    GROUP BY source, kind
	) sub ON r.source = sub.source AND r.kind = sub.kind AND r.version = sub.latest_version
	GROUP BY r.source, r.kind, m.key, m.value
	HAVING COUNT(distinct m.resource_id) > 2
	ORDER BY r.source, r.kind, m.key, total_count DESC;
	`

	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MetadataCount
	for rows.Next() {
		var mkvc MetadataCount
		err := rows.Scan(&mkvc.Source, &mkvc.Kind, &mkvc.Key, &mkvc.Value, &mkvc.TotalCount)
		if err != nil {
			return nil, err
		}
		results = append(results, mkvc)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
