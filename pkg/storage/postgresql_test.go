package storage

import (
	"errors"
	"reflect"
	"testing"

	"github.com/shahariaazam/teredix/pkg/config"
	"github.com/shahariaazam/teredix/pkg/resource"

	"github.com/stretchr/testify/assert"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestPostgreSQL_Prepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	storage := &PostgreSQL{DB: db}

	// Expect the resources table to be created
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS resources(.*)CREATE TABLE IF NOT EXISTS metadata(.*)CREATE TABLE IF NOT EXISTS relations`).WillReturnResult(sqlmock.NewResult(0, 0))

	if err := storage.Prepare(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestPostgreSQL_Persist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %s", err)
	}
	defer db.Close()

	pg := &PostgreSQL{DB: db}

	mock.ExpectBegin()

	// mock Persist statements
	resourcesStmt := mock.ExpectPrepare(`INSERT INTO resources`)
	metadataStmt := mock.ExpectPrepare(`INSERT INTO metadata`)

	resources := []resource.Resource{
		{
			Kind:       "resource1",
			UUID:       "uuid1",
			Name:       "name1",
			ExternalID: "external_id1",
			MetaData: []resource.MetaData{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
			RelatedWith: []resource.Resource{
				{
					Kind:       "resource2",
					UUID:       "uuid2",
					Name:       "name2",
					ExternalID: "external_id2",
				},
			},
		},
	}

	// mock Persist statement results
	resourcesStmt.ExpectQuery().WithArgs("resource1", "uuid1", "name1", "external_id1").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	metadataStmt.ExpectExec().WithArgs(1, "key1", "value1").WillReturnResult(sqlmock.NewResult(0, 0))

	// call the method being tested
	err = pg.Persist(resources)
	if err != nil {
		t.Errorf("error calling Persist: %s", err)
	}

	// verify that the mock expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %s", err)
	}
}

func TestPostgreSQL_Find(t *testing.T) {
	testCases := []struct {
		name                  string
		resourceFilter        ResourceFilter
		expectedQuery         string
		expectedResourceCount int
	}{
		{
			name:                  "find without any filtering parameter",
			resourceFilter:        ResourceFilter{},
			expectedQuery:         "SELECT r.kind, r.uuid, r.name, r.external_id, m.key, m.value, rr.kind, rr.uuid, rr.name, rr.external_id FROM resources r LEFT JOIN metadata m ON r.id = m.resource_id LEFT JOIN relations rl ON r.id = rl.resource_id LEFT JOIN resources rr ON rl.related_resource_id = rr.id",
			expectedResourceCount: 2,
		},
		{
			name:                  "find by resource name",
			resourceFilter:        ResourceFilter{Name: "name1"},
			expectedQuery:         `SELECT r.kind, r.uuid, r.name, r.external_id, m.key, m.value, rr.kind, rr.uuid, rr.name, rr.external_id FROM resources r LEFT JOIN metadata m ON r.id = m.resource_id LEFT JOIN relations rl ON r.id = rl.resource_id LEFT JOIN resources rr ON rl.related_resource_id = rr.id WHERE r.name = \$1`,
			expectedResourceCount: 0,
		},
		{
			name:                  "find by resource uuid",
			resourceFilter:        ResourceFilter{UUID: "uuid1"},
			expectedQuery:         `SELECT r.kind, r.uuid, r.name, r.external_id, m.key, m.value, rr.kind, rr.uuid, rr.name, rr.external_id FROM resources r LEFT JOIN metadata m ON r.id = m.resource_id LEFT JOIN relations rl ON r.id = rl.resource_id LEFT JOIN resources rr ON rl.related_resource_id = rr.id WHERE r.uuid = \$1`,
			expectedResourceCount: 0,
		},
		{
			name:                  "find by resource ExternalID",
			resourceFilter:        ResourceFilter{ExternalID: "external_id1"},
			expectedQuery:         `SELECT r.kind, r.uuid, r.name, r.external_id, m.key, m.value, rr.kind, rr.uuid, rr.name, rr.external_id FROM resources r LEFT JOIN metadata m ON r.id = m.resource_id LEFT JOIN relations rl ON r.id = rl.resource_id LEFT JOIN resources rr ON rl.related_resource_id = rr.id WHERE r.external_id`,
			expectedResourceCount: 0,
		},
		{
			name:                  "find by resource kind, uuid, name",
			resourceFilter:        ResourceFilter{Kind: "kind1", UUID: "uuid1", Name: "name1"},
			expectedQuery:         `SELECT r.kind, r.uuid, r.name, r.external_id, m.key, m.value, rr.kind, rr.uuid, rr.name, rr.external_id FROM resources r LEFT JOIN metadata m ON r.id = m.resource_id LEFT JOIN relations rl ON r.id = rl.resource_id LEFT JOIN resources rr ON rl.related_resource_id = rr.id WHERE r\.kind = .+? AND r\.uuid = .+? AND r\.name = `,
			expectedResourceCount: 0,
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	p := &PostgreSQL{DB: db}

	// expected rows returned by the mock database query
	expectedRows := sqlmock.NewRows([]string{"kind", "uuid", "name", "external_id", "meta_key", "meta_value", "related_kind", "related_uuid", "related_name", "related_external_id"}).
		AddRow("kind1", "uuid1", "name1", "external_id1", "meta_key1", "meta_value1", "related_kind1", "related_uuid1", "related_name1", "related_external_id1").
		AddRow("kind2", "uuid2", "name2", "external_id2", "meta_key2", "meta_value2", "related_kind2", "related_uuid2", "related_name2", "related_external_id2")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set up mock database query expectations
			mock.ExpectQuery(tc.expectedQuery).
				WillReturnRows(expectedRows)

			// call the method being tested
			resources, err := p.Find(tc.resourceFilter)
			if err != nil {
				t.Fatalf("unexpected error from Find: %v", err)
			}

			// verify the result
			if len(resources) != tc.expectedResourceCount {
				t.Fatalf("unexpected number of resources: got %d, want %d", len(resources), tc.expectedResourceCount)
			}
		})
	}
}

func TestPostgreSQL_Find_For_Error(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	p := &PostgreSQL{DB: db}

	mock.ExpectQuery("SELECT r.kind, r.uuid, r.name, r.external_id, m.key, m.value, rr.kind, rr.uuid, rr.name, rr.external_id FROM resources r LEFT JOIN metadata m ON r.id = m.resource_id LEFT JOIN relations rl ON r.id = rl.resource_id LEFT JOIN resources rr ON rl.related_resource_id = rr.id").
		WillReturnError(errors.New("query failed"))

	// call the method being tested
	resources, err := p.Find(ResourceFilter{})
	assert.Error(t, err)

	assert.Equal(t, 0, len(resources))
}

func TestGetRelations(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database connection: %s", err)
	}
	defer db.Close()

	// Create a new PostgreSQL object using the mock database connection
	p := &PostgreSQL{DB: db}

	// Set up expected results from the mock database
	expected := []map[string]string{
		{"resource1": "related_resource1"},
		{"resource2": "related_resource2"},
	}

	// Set up mock database rows and columns
	rows := sqlmock.NewRows([]string{"resource_uuid", "related_resource_uuid"}).
		AddRow("resource1", "related_resource1").
		AddRow("resource2", "related_resource2")

	// Expect a query to the "relations" table and return the mock rows
	mock.ExpectQuery("SELECT r.uuid as resource_uuid, r2.uuid as related_resource_uuid FROM relations left join resources r on r.id = relations.resource_id left join resources r2 on r2.id = relations.related_resource_id").
		WillReturnRows(rows)

	// Call the GetRelations method
	relations, err := p.GetRelations()

	// Check for errors
	if err != nil {
		t.Fatalf("Error calling GetRelations: %s", err)
	}

	// Compare expected and actual results
	if !reflect.DeepEqual(expected, relations) {
		t.Fatalf("Expected %v, but got %v", expected, relations)
	}

	// Verify that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Unfulfilled expectations: %s", err)
	}
}

func TestAnalyzeRelationMatrix(t *testing.T) {
	// Create a new mock database and get a handle for the mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a new instance of the PostgreSQL struct with the mock database
	p := &PostgreSQL{DB: db}

	// Define expected rows for the first query
	rows1 := sqlmock.NewRows([]string{"resource_ids"}).
		AddRow("1,2,3")

	// Define expected rows for the second query
	rows2 := sqlmock.NewRows([]string{"resource_ids"}).
		AddRow("4,5,6")

	// Set up the mock expectations for the first query
	mock.ExpectQuery("SELECT STRING_AGG").
		WithArgs("related_key", "related_value", "related").
		WillReturnRows(rows1)

	// Set up the mock expectations for the second query
	mock.ExpectQuery("SELECT STRING_AGG").
		WithArgs("metadata_key", "metadata_value", "kind").
		WillReturnRows(rows2)

	// Call the function with the mock PostgreSQL instance
	matrix, err := p.analyzeRelationMatrix(config.RelationCriteria{
		RelatedMetadataKey:   "related_key",
		RelatedMetadataValue: "related_value",
		RelatedKind:          "related",
		MetadataKey:          "metadata_key",
		MetadataValue:        "metadata_value",
		Kind:                 "kind",
		Source: config.RelationCriteriaNode{
			Kind:      "kind",
			MetaKey:   "metadata_key",
			MetaValue: "metadata_value",
		},
		Target: config.RelationCriteriaNode{
			Kind:      "related",
			MetaKey:   "related_key",
			MetaValue: "related_value",
		},
	})

	// Check the results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(matrix) != 9 {
		t.Errorf("Unexpected number of rows returned: %d", len(matrix))
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
