//go:build integration

package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const yamlContent = `
---
organization:
  name: Shaharia Lab
  logo: http://example.com

discovery:
  name: Infrastructure Discovery
  description: Some description text
  worker_pool_size: 1

storage:
  batch_size: 1
  engines:
    postgresql:
      host: "localhost"
      port: 5432
      user: "app"
      password: "pass"
      db: "app"
    neo4j:
      config_key: "value"
  default_engine: postgresql

source:
  fs_one:
    type: file_system
    configuration:
      root_directory: "{{ROOT_DIRECTORY}}"
    fields: &file_system_fields
      - rootDirectory
      - machineHost
    schedule: "@every 10s"
  fs_two:
    type: file_system
    configuration:
      root_directory: "{{ROOT_DIRECTORY}}"
    fields: &file_system_fields
      - rootDirectory
      - machineHost
    schedule: "@every 10s"
relations:
  criteria:
    - name: "file-system-rule1"
      source:
        kind: "FilePath"
        meta_key: "rootDirectory"
        meta_value: "/some/path"
      target:
        kind: "FilePath"
        meta_key: "rootDirectory"
        meta_value: "/some/path"
`

func Test_run(t *testing.T) {
	commonAppConfig := &config.AppConfig{
		Organization: config.Organization{
			Name: "Shaharia Lab",
			Logo: "http://example.com",
		},
		Discovery: config.Discovery{
			Name:           "Infrastructure Discovery",
			Description:    "Some description text",
			WorkerPoolSize: 1,
		},
		Storage: config.Storage{
			BatchSize: 1,
			Engines: map[string]interface{}{
				"postgresql": map[string]interface{}{
					"host":     "localhost",
					"port":     5432,
					"user":     "app",
					"password": "pass",
					"db":       "app",
				},
			},
			DefaultEngine: "postgresql",
		},
		Sources: map[string]config.Source{
			"fs_one": {
				Type: "file_system",
				Configuration: map[string]string{
					"root_directory": "",
				},
				Fields:   []string{"rootDirectory", "machineHost"},
				Schedule: "@every 2s",
			},
		},
		Relation: config.Relation{
			RelationCriteria: []config.RelationCriteria{
				{
					Name: "file-system-rule1",
					Source: config.RelationCriteriaNode{
						Kind:      "FilePath",
						MetaKey:   "rootDirectory",
						MetaValue: "/some/path",
					},
					Target: config.RelationCriteriaNode{
						Kind:      "FilePath",
						MetaKey:   "rootDirectory",
						MetaValue: "/some/path",
					},
				},
			},
		},
	}

	testCases := []struct {
		name         string
		sourceConfig map[string]config.Source
	}{
		{
			name: "valid configuration",
			sourceConfig: map[string]config.Source{
				"fs_one": {
					Type: "file_system",
					Configuration: map[string]string{
						"root_directory": getTestRootDiretory("cmd/testdata"),
					},
					Fields:   []string{"rootDirectory", "machineHost"},
					Schedule: "@every 2s",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testDBHost := os.Getenv("TEST_DB_HOST")
			if testDBHost == "" {
				testDBHost = "localhost"
			}

			commonAppConfig.Storage.Engines["postgresql"].(map[string]interface{})["host"] = testDBHost

			commonAppConfig.Sources = tc.sourceConfig

			configFilePath, err := SaveToTempYAML(commonAppConfig)
			if err != nil {
				return
			}

			appConfig, err := config.Load(configFilePath)
			assert.NoError(t, err)

			// Create a context with a timeout.
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel() // Always defer the cancellation, which ensures resources are cleaned up.

			err = run(ctx, appConfig, &logrus.Logger{})
			assert.NoError(t, err)

			defer resetDatabase(testDBHost)
		})
	}
}

func getTestRootDiretory(directory string) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	return filepath.Join(dir, directory)
}

func SaveToTempYAML(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data to YAML: %w", err)
	}

	// Create a temporary file
	tmpfile, err := ioutil.TempFile("", "data-*.yaml")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	// Close the file after writing
	defer tmpfile.Close()

	if _, err := tmpfile.Write(bytes); err != nil {
		os.Remove(tmpfile.Name()) // Clean up in case of an error
		return "", fmt.Errorf("failed to write to temporary file: %w", err)
	}

	return tmpfile.Name(), nil
}

func resetDatabase(testDBHost string) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		testDBHost,
		5432,
		"app",
		"pass",
		"app",
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	deleteTables(db, []string{"metadata", "relations", "resources"})
}

func deleteTables(db *sql.DB, tables []string) error {
	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", table)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error deleting table %s: %v", table, err)
		}
		fmt.Printf("Deleted table %s successfully\n", table)
	}
	return nil
}
