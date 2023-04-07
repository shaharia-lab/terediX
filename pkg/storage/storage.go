// Package storage store resource information
package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/shahariaazam/teredix/pkg/config"
	"github.com/shahariaazam/teredix/pkg/resource"
)

// ResourceFilter configure the filter
type ResourceFilter struct {
	Kind       string
	UUID       string
	Name       string
	ExternalID string
}

// Storage interface helps to build different storage
type Storage interface {
	// Prepare to prepare the Storage schema
	Prepare() error

	// Persist Save resources to Storage
	Persist(resources []resource.Resource) error

	// Find will return the resources based on ResourceFilter
	Find(filter ResourceFilter) ([]resource.Resource, error)

	GetResources() ([]resource.Resource, error)
	GetRelations() ([]map[string]string, error)

	// StoreRelations Store Relationship
	StoreRelations(relation config.Relation) error
}

// Query build query based on filters
type Query struct {
	filters []string
	params  []interface{}
}

// AddFilter adds filter
func (q *Query) AddFilter(field, operator string, value interface{}) {
	q.filters = append(q.filters, fmt.Sprintf("%s %s $%d", field, operator, len(q.params)+1))
	q.params = append(q.params, value)
}

// Build builds query
func (q *Query) Build() (string, []interface{}) {
	var whereClause string
	if len(q.filters) > 0 {
		whereClause = "WHERE " + strings.Join(q.filters, " AND ")
	}
	return fmt.Sprintf("SELECT r.kind, r.uuid, r.name, r.external_id, m.key, m.value, rr.kind, rr.uuid, rr.name, rr.external_id FROM resources r LEFT JOIN metadata m ON r.id = m.resource_id LEFT JOIN relations rl ON r.id = rl.resource_id LEFT JOIN resources rr ON rl.related_resource_id = rr.id %s", whereClause), q.params
}

// BuildStorage build storage based on configuration
func BuildStorage(appConfig *config.AppConfig) Storage {
	var st Storage

	for engineKey, engine := range appConfig.Storage.Engines {
		if engineKey == "postgresql" {
			ec, _ := engine.(map[string]interface{})
			connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				ec["host"],
				ec["port"],
				ec["user"],
				ec["password"],
				ec["db"],
			)

			db, err := sql.Open("postgres", connStr)
			if err != nil {
				panic(ec)
			}

			st = &PostgreSQL{DB: db}
		}
	}

	return st
}
