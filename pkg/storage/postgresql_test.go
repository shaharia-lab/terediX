package storage

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"teredix/pkg/resource"
	"testing"
)

func TestPostgreSQL_Prepare(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	storage := &PostgreSQL{DB: db}

	// Expect the resources table to be created
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS resources`).WillReturnResult(sqlmock.NewResult(0, 0))
	// Expect the metadata table to be created
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS metadata`).WillReturnResult(sqlmock.NewResult(0, 0))
	// Expect the relations table to be created
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS relations`).WillReturnResult(sqlmock.NewResult(0, 0))

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
	pg.Prepare()

	mock.ExpectBegin()

	// mock Persist statements
	resourcesStmt := mock.ExpectPrepare(`INSERT INTO resources`)
	metadataStmt := mock.ExpectPrepare(`INSERT INTO metadata`)
	relationsStmt := mock.ExpectPrepare(`INSERT INTO relations`)

	// define resources to be persisted
	//now := time.Now()
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
	relationsStmt.ExpectExec().WithArgs(1, "external_id2").WillReturnResult(sqlmock.NewResult(0, 0))

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
