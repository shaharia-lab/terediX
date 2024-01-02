package storage

import (
	"reflect"
	"testing"

	"github.com/shaharia-lab/teredix/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestAddFilter(t *testing.T) {
	q := &Query{}

	q.AddFilter("kind", "=", "vm")
	q.AddFilter("name", "LIKE", "%web%")

	expectedQuery := `SELECT  r.source, r.kind, r.uuid, r.name, r.external_id, r.version, r.discovered_at, json_object_agg(m.key, m.value) FILTER (WHERE m.key IS NOT NULL) AS meta_data
FROM 
    resources r
LEFT JOIN 
    metadata m ON r.id = m.resource_id
WHERE kind = $1 AND name LIKE $2
GROUP BY r.id LIMIT 200 OFFSET 0
`
	expectedParams := []interface{}{"vm", "%web%"}

	queryString, params := q.Build()

	assert.Equal(t, expectedQuery, queryString)
	assert.Equal(t, expectedParams, params)
}

func TestBuild(t *testing.T) {
	q := &Query{}

	// Test an empty query
	query, params := q.Build()
	assert.Equal(t, `SELECT  r.source, r.kind, r.uuid, r.name, r.external_id, r.version, r.discovered_at, json_object_agg(m.key, m.value) FILTER (WHERE m.key IS NOT NULL) AS meta_data
FROM 
    resources r
LEFT JOIN 
    metadata m ON r.id = m.resource_id

GROUP BY r.id LIMIT 200 OFFSET 0
`, query)
	assert.Empty(t, params)

	// Test a query with one filter
	q.AddFilter("kind", "=", "test")
	query, params = q.Build()
	assert.Equal(t, `SELECT  r.source, r.kind, r.uuid, r.name, r.external_id, r.version, r.discovered_at, json_object_agg(m.key, m.value) FILTER (WHERE m.key IS NOT NULL) AS meta_data
FROM 
    resources r
LEFT JOIN 
    metadata m ON r.id = m.resource_id
WHERE kind = $1
GROUP BY r.id LIMIT 200 OFFSET 0
`, query)
	assert.Equal(t, []interface{}{"test"}, params)
	// Test a query with multiple filters
	q.AddFilter("name", "=", "test-resource")
	q.AddFilter("uuid", "!=", "1234")
	q.AddFilter("external_id", "like", "%abc%")
	query, params = q.Build()
	assert.Equal(t, `SELECT  r.source, r.kind, r.uuid, r.name, r.external_id, r.version, r.discovered_at, json_object_agg(m.key, m.value) FILTER (WHERE m.key IS NOT NULL) AS meta_data
FROM 
    resources r
LEFT JOIN 
    metadata m ON r.id = m.resource_id
WHERE kind = $1 AND name = $2 AND uuid != $3 AND external_id like $4
GROUP BY r.id LIMIT 200 OFFSET 0
`, query)
	assert.Equal(t, []interface{}{"test", "test-resource", "1234", "%abc%"}, params)
}

func TestBuildStorage(t *testing.T) {
	// Set up test configuration
	appConfig := &config.AppConfig{
		Storage: config.Storage{
			Engines: map[string]interface{}{
				"postgresql": map[string]interface{}{
					"host":     "localhost",
					"port":     5432,
					"user":     "testuser",
					"password": "testpassword",
					"db":       "testdb",
				},
			},
		},
	}

	// Call BuildStorage
	st, _ := BuildStorage(appConfig)

	// Verify that the correct type of storage was created
	expectedType := reflect.TypeOf(&PostgreSQL{})
	actualType := reflect.TypeOf(st)
	if expectedType != actualType {
		t.Errorf("Expected type %v but got %v", expectedType, actualType)
	}

	// Verify that the storage was properly initialized
	_, ok := st.(*PostgreSQL)
	if !ok {
		t.Errorf("Expected storage to be of type PostgreSQL")
	}
}
