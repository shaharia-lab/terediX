package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test with valid config file
	path := "testdata/valid_config.yaml"
	_, err := Load(path)
	assert.NoError(t, err)

	// Test with non-existent config file
	path = "testdata/non_existent_config.yaml"
	_, err = Load(path)
	assert.Error(t, err)

	// Test with invalid YAML data
	path = "testdata/invalid_config.yaml"
	_, err = Load(path)
	assert.Error(t, err)
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name    string
		config  AppConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing organization name",
			config: AppConfig{
				Discovery: Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing discovery name",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing storage engine",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "non_existent_engine",
					Engines:       map[string]interface{}{},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},

			wantErr: true,
		},
		{
			name: "missing storage default engine",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "non_existent_engine",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid storage engine config",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid default storage engine",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "invalid_engine",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid batch_size",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     -1,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing source type",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing kube_config_file_path for kubernetes config",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "kubernetes",
						Configuration: map[string]string{
							"invalid_key": "invalid_value",
						},
					},
					"source2": {
						Type: "kubernetes",
						Configuration: map[string]string{
							"kube_config_file_path": "",
						},
					},
					"source3": {
						Type: "kubernetes",
						Configuration: map[string]string{
							"kube_config_file_path": "/path/to/kube/config",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing kube_config_file_path for kubernetes config",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source3": {
						Type: "kubernetes",
						Configuration: map[string]string{
							"kube_config_file_path": "/path/to/kube/config",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing root_directory for file_system source",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type:          "file_system",
						Configuration: map[string]string{},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty root_directory for file_system source",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing neo4j engine config_key",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "neo4j",
					Engines: map[string]interface{}{
						"neo4j": map[string]interface{}{
							"host":     "localhost",
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid neo4j engine config",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "neo4j",
					Engines: map[string]interface{}{
						"neo4j": map[string]interface{}{
							"config_key": "myconfig",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid neo4j engine configuration: not a map",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "neo4j",
					Engines: map[string]interface{}{
						"neo4j": "invalid_configuration",
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid postgresql engine configuration - not a map",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": []string{"invalid config"},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "unknown storage engine",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "cassandra", // this engine is not defined
					Engines: map[string]interface{}{
						"cassandra": map[string]interface{}{
							"db": "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no sources defined",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "unknown source type",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "unknown_type",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid source depends_on",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
						DependsOn: []string{"source2"},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid source depends_on - multiple dependencies",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
						DependsOn: []string{"source2", "source3"},
					},
				},
				RelationCriteria: []RelationCriteria{
					{
						Name:                 "name",
						Kind:                 "kind",
						MetadataKey:          "source-kind-key1",
						MetadataValue:        "source-kind-value1",
						RelatedKind:          "related-kind",
						RelatedMetadataKey:   "related-metadata-key",
						RelatedMetadataValue: "related-metadata-value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing relations field",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing relations field",
			config: AppConfig{
				Organization: Organization{Name: "My Org", Logo: "http://example.com"},
				Discovery:    Discovery{Name: "My Discovery", Description: "Discovery description"},
				Storage: Storage{
					BatchSize:     2,
					DefaultEngine: "postgresql",
					Engines: map[string]interface{}{
						"postgresql": map[string]interface{}{
							"host":     "localhost",
							"port":     5432,
							"user":     "myuser",
							"password": "mypassword",
							"db":       "mydb",
						},
					},
				},
				Sources: map[string]Source{
					"source1": {
						Type: "file_system",
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(&tc.config)
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
