// Package storage store resource information
package storage

import (
	"database/sql"
	"strings"

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
        kind TEXT NOT NULL,
        uuid TEXT NOT NULL,
        name TEXT NOT NULL,
        external_id TEXT NOT NULL UNIQUE,
        discovered_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS metadata (
        resource_id INTEGER REFERENCES resources(id),
        key TEXT NOT NULL,
        value TEXT NOT NULL,
        PRIMARY KEY(resource_id, key, value)
    );

    CREATE TABLE IF NOT EXISTS relations (
        resource_id INTEGER REFERENCES resources(id),
        related_resource_id INTEGER REFERENCES resources(id),
        PRIMARY KEY(resource_id, related_resource_id)
    );

	CREATE INDEX IF NOT EXISTS metadata_key_idx ON metadata (key);
	CREATE INDEX IF NOT EXISTS metadata_value_idx ON metadata (value);
`

	_, err := p.DB.Exec(sqlString)
	if err != nil {
		return err
	}

	return nil
}

// Persist store resources
func (p *PostgreSQL) Persist(resources []resource.Resource) error {
	return p.runInTransaction(func(tx *sql.Tx) error {
		// Prepare the SQL statements
		resourcesStmt, err := tx.Prepare("INSERT INTO resources (kind, uuid, name, external_id) VALUES ($1, $2, $3, $4) ON CONFLICT (external_id) DO UPDATE SET kind = excluded.kind, uuid = excluded.uuid, name = excluded.name RETURNING id")
		if err != nil {
			return err
		}
		defer func(resourcesStmt *sql.Stmt) {
			err := resourcesStmt.Close()
			if err != nil {
				return
			}
		}(resourcesStmt)

		metadataStmt, err := tx.Prepare("INSERT INTO metadata (resource_id, key, value) VALUES ($1, $2, $3) ON CONFLICT (resource_id, key, value) DO UPDATE SET value = excluded.value")
		if err != nil {
			return err
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
			err := resourcesStmt.QueryRow(res.Kind, res.UUID, res.Name, res.ExternalID).Scan(&id)
			if err != nil {
				return err
			}

			// Insert or update the metadata
			for _, meta := range res.MetaData {
				_, err = metadataStmt.Exec(id, meta.Key, meta.Value)
				if err != nil {
					return err
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
		return resources, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	// Parse the results
	for rows.Next() {
		var kind, uuid, name, externalID, metaKey, metaValue string
		var relatedKind, relatedUUID, relatedName, relatedExternalID sql.NullString

		err := rows.Scan(&kind, &uuid, &name, &externalID, &metaKey, &metaValue, &relatedKind, &relatedUUID, &relatedName, &relatedExternalID)
		if err != nil {
			return resources, err
		}

		// Create a resource object if it doesn't exist in the slice yet
		var res *resource.Resource
		for i := range resources {
			if resources[i].UUID == uuid {
				res = &resources[i]
				break
			}
		}
		if res == nil {
			res = &resource.Resource{
				Kind:        kind,
				UUID:        uuid,
				Name:        name,
				ExternalID:  externalID,
				MetaData:    []resource.MetaData{},
				RelatedWith: []resource.Resource{},
			}
			resources = append(resources, *res)
		}

		// Add metadata to the resource
		if metaKey != "" && metaValue != "" {
			res.MetaData = append(res.MetaData, resource.MetaData{
				Key:   metaKey,
				Value: metaValue,
			})
		}

		// Add related resource to the resource
		if relatedKind.Valid && relatedKind.String != "" && relatedUUID.String != "" {
			related := resource.Resource{
				Kind:       relatedKind.String,
				UUID:       relatedUUID.String,
				Name:       relatedName.String,
				ExternalID: relatedExternalID.String,
			}
			res.RelatedWith = append(res.RelatedWith, related)
		}
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
			return err
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
				return err
			}
		}

		return nil
	})
}

func (p *PostgreSQL) storeRelationMatrix(rc config.RelationCriteria, relationsStmt *sql.Stmt) error {
	relationMatrix, err := p.analyzeRelationMatrix(rc)
	if err != nil {
		return err
	}

	// Insert relationships for the matching resources
	for _, c := range relationMatrix {
		for k, v := range c {
			if _, err := relationsStmt.Exec(k, v); err != nil {
				return err
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
	resources := []resource.Resource{}

	// Query all resources
	rows, err := p.DB.Query("SELECT kind, uuid, name, external_id FROM resources")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	// Loop through the result set and create Resource objects
	for rows.Next() {
		r := resource.Resource{}
		if err := rows.Scan(&r.Kind, &r.UUID, &r.Name, &r.ExternalID); err != nil {
			return nil, err
		}
		resources = append(resources, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}

// GetRelations fetch all the relations between resources
func (p *PostgreSQL) GetRelations() ([]map[string]string, error) {
	var relations []map[string]string

	// Query all relations
	rows, err := p.DB.Query("SELECT r.uuid as resource_uuid, r2.uuid as related_resource_uuid FROM relations left join resources r on r.id = relations.resource_id left join resources r2 on r2.id = relations.related_resource_id")
	if err != nil {
		return nil, err
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
			return nil, err
		}
		r[resourceUUID] = relatedResourceUUID
		relations = append(relations, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return relations, nil
}

func (p *PostgreSQL) runInTransaction(f func(tx *sql.Tx) error) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
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
