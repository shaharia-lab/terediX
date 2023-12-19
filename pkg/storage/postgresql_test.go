package storage

import (
	"reflect"
	"testing"

	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"

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

	r1 := resource.NewResource("resource1", "name1", "external_id1", "scanner_name", 1)
	r1.SetUUID("uuid1")
	r1.AddMetaData(map[string]string{"key1": "value1"})
	r1.AddRelation(resource.NewResource("resource2", "name2", "external_id2", "", 1))

	resources := []resource.Resource{r1}

	// mock Persist statement results
	resourcesStmt.ExpectQuery().WithArgs("resource1", "uuid1", "name1", "external_id1", "scanner_name", 1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
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
