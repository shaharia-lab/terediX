//go:build integration

package processor

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/scanner"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/stretchr/testify/assert"
)

const fileSystemSourceIntegrationTestRootDirectory = "config/testdata"

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

func TestProcessor_Process_Integration(t *testing.T) {
	testDBHost := os.Getenv("TEST_DB_HOST")
	if testDBHost == "" {
		testDBHost = "localhost"
	}

	yamlContentReplaced := strings.ReplaceAll(yamlContent, `host: "localhost"`, fmt.Sprintf(`host: "%s"`, testDBHost))
	yamlContentReplaced = strings.ReplaceAll(yamlContentReplaced, `{{ROOT_DIRECTORY}}`, getTestRootDiretory(t))

	WriteToFile("config.yaml", yamlContentReplaced)
	appConfig, err := config.Load("config.yaml")
	assert.NoError(t, err)

	sources := scanner.NewSourceRegistry(scanner.GetScannerRegistries())

	st := storage.BuildStorage(appConfig)
	err = st.Prepare()
	assert.NoError(t, err)

	processConfig := Config{BatchSize: appConfig.Storage.BatchSize}
	p := NewProcessor(processConfig, st, sources.BuildFromAppConfig(appConfig.Sources))

	resourceChan := make(chan resource.Resource)
	p.Process(resourceChan)

	time.Sleep(2 * time.Second)

	resources, err := st.Find(storage.ResourceFilter{Kind: "FilePath"})
	assert.NoError(t, err)

	assert.Equal(t, 2, len(resources))

	resetDatabase(testDBHost)
}

// WriteToFile writes content to a given filename.
func WriteToFile(filename string, content string) error {
	// Convert the content string to a byte slice as required by ioutil.WriteFile
	data := []byte(content)

	// Write data to the file
	// This will create the file if it doesn't exist, or truncate the file if it does
	err := ioutil.WriteFile(filename, data, 0644)
	return err
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

func getTestRootDiretory(t *testing.T) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	return filepath.Join(dir, fileSystemSourceIntegrationTestRootDirectory)
}
