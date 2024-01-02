// Package storage store resource information
package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"
)

const noOfVersionToKeep = 2

// ResourceFilter configure the filter
type ResourceFilter struct {
	Kind       string
	UUID       string
	Name       string
	ExternalID string
	PerPage    int
	Offset     int
}

// ResourceCount count resource
type ResourceCount struct {
	Source     string
	Kind       string
	TotalCount int
}

// MetadataCount count resources by metadata
type MetadataCount struct {
	Source     string
	Kind       string
	Key        string
	Value      string
	TotalCount int
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

	GetNextVersionForResource(source, kind string) (int, error)

	CleanupOldVersion(source, kind string) (int64, error)

	GetResourceCount() ([]ResourceCount, error)

	GetResourceCountByMetaData() ([]MetadataCount, error)
}

// Query build query based on filters
type Query struct {
	filters []string
	params  []interface{}
	perPage int
	offset  int
}

// AddFilter adds filter
func (q *Query) AddFilter(field, operator string, value interface{}) {
	q.filters = append(q.filters, fmt.Sprintf("%s %s $%d", field, operator, len(q.params)+1))
	q.params = append(q.params, value)
}

// SetPerPage set per page
func (q *Query) SetPerPage(perPage int) {
	q.perPage = perPage
}

// SetOffset set offset
func (q *Query) SetOffset(from int) {
	q.offset = from
}

// Build builds query
func (q *Query) Build() (string, []interface{}) {

	if q.perPage == 0 {
		q.perPage = 200
	}

	if q.offset == 0 {
		q.offset = 0
	}

	var whereClause string
	if len(q.filters) > 0 {
		whereClause = "WHERE " + strings.Join(q.filters, " AND ")
	}

	query := fmt.Sprintf(`SELECT  r.source, r.kind, r.uuid, r.name, r.external_id, r.version, r.discovered_at, json_object_agg(m.key, m.value) FILTER (WHERE m.key IS NOT NULL) AS meta_data
FROM 
    resources r
LEFT JOIN 
    metadata m ON r.id = m.resource_id
%s
GROUP BY r.id LIMIT %d OFFSET %d
`, whereClause, q.perPage, q.offset)

	return query, q.params
}

// BuildStorage build storage based on configuration
func BuildStorage(appConfig *config.AppConfig) (Storage, error) {
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
				return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
			}

			st = &PostgreSQL{DB: db}
		}
	}

	return st, nil
}
