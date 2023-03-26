//pkg/storage/postgresql/postgresql.go
package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"teredix/pkg/config"
	"teredix/pkg/resource"
	"time"
)

type PostgreSQL struct {
	DB *sql.DB
}

func (p *PostgreSQL) Prepare() error {
	// Create the resources table if it doesn't exist
	_, err := p.DB.Exec(`CREATE TABLE IF NOT EXISTS resources (
		id SERIAL PRIMARY KEY,
		kind TEXT NOT NULL,
		uuid TEXT NOT NULL,
		name TEXT NOT NULL,
		external_id TEXT NOT NULL UNIQUE,
		discovered_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return err
	}

	// Create the metadata table if it doesn't exist
	_, err = p.DB.Exec(`CREATE TABLE IF NOT EXISTS metadata (
		resource_id INTEGER REFERENCES resources(id),
		key TEXT NOT NULL,
		value TEXT NOT NULL,
		PRIMARY KEY(resource_id, key)
	)`)
	if err != nil {
		return err
	}

	// Create the relations table if it doesn't exist
	_, err = p.DB.Exec(`CREATE TABLE IF NOT EXISTS relations (
		resource_id INTEGER REFERENCES resources(id),
		related_resource_id INTEGER REFERENCES resources(id),
		PRIMARY KEY(resource_id, related_resource_id)
	)`)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) Persist(resources []resource.Resource) error {
	// Begin a transaction
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

	// Prepare the SQL statements
	resourcesStmt, err := tx.Prepare("INSERT INTO resources (kind, uuid, name, external_id, discovered_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (external_id) DO UPDATE SET kind = excluded.kind, uuid = excluded.uuid, name = excluded.name RETURNING id")
	if err != nil {
		return err
	}
	defer resourcesStmt.Close()

	metadataStmt, err := tx.Prepare("INSERT INTO metadata (resource_id, key, value) VALUES ($1, $2, $3) ON CONFLICT (resource_id, key) DO UPDATE SET value = excluded.value")
	if err != nil {
		return err
	}
	defer metadataStmt.Close()

	relationsStmt, err := tx.Prepare("INSERT INTO relations (resource_id, related_resource_id) SELECT $1, r.id FROM resources r WHERE r.external_id = $2 ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}
	defer relationsStmt.Close()

	// Loop through the resources and insert or update them into the database
	for _, res := range resources {
		// Insert or update the resource
		var id int
		err := resourcesStmt.QueryRow(res.Kind, res.UUID, res.Name, res.ExternalID, time.Now()).Scan(&id)
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

		// Insert or update the relations
		for _, related := range res.RelatedWith {
			_, err = relationsStmt.Exec(id, related.ExternalID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

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
	defer rows.Close()

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
func (p *PostgreSQL) StoreRelations(relation config.Relation) error {
	// Begin a transaction
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

	// Prepare the SQL statement to select resources based on the relation
	selectStmt := `SELECT r1.id AS id1, r2.id AS id2 FROM resources r1, resources r2
		WHERE r1.kind = $1 AND r1.external_id IN (
			SELECT resource_id FROM metadata WHERE key = $2 AND value = $3
		) AND r2.kind = $4 AND r2.external_id IN (
			SELECT resource_id FROM metadata WHERE key = $5 AND value = $6
		) AND r1.id != r2.id`

	// Prepare the SQL statement to insert relationships into the relations table
	relationsStmt, err := tx.Prepare("INSERT INTO relations (resource_id, related_resource_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer relationsStmt.Close()

	// Loop through each relation and find the resources that match it
	for _, c := range relation.RelationCriteria {
		rows, err := p.DB.Query(selectStmt, c.Kind, c.MetadataKey, c.MetadataValue, c.RelatedKind, c.RelatedMetadataKey, c.RelatedMetadataValue)
		if err != nil {
			return err
		}

		// Insert relationships for the matching resources
		for rows.Next() {
			var id1, id2 int
			if err := rows.Scan(&id1, &id2); err != nil {
				return err
			}
			if _, err := relationsStmt.Exec(id1, id2); err != nil {
				return err
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}
	}

	return nil
}
