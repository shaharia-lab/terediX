package config

import (
	"testing"

	"github.com/shahariaazam/teredix/pkg"

	"github.com/stretchr/testify/assert"
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

func TestLoadAndValidateConfigYaml(t *testing.T) {
	// Test with valid config file
	path := "testdata/valid_config.yaml"
	appConfig, loadError := Load(path)
	assert.NoError(t, loadError)

	validationError := Validate(appConfig)
	assert.NoError(t, validationError)
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
					"source2": {
						Type: pkg.SourceTypeGitHubRepository,
						Configuration: map[string]string{
							"token":         "mytoken",
							"user_or_org":   "myuser",
							"repository":    "myrepo",
							"branch":        "mybranch",
							"path":          "mypath",
							"file_patterns": "*.yaml",
						},
					},
					"source3": {
						Type: pkg.SourceTypeAWSS3,
						Configuration: map[string]string{
							"access_key":    "access_key",
							"secret_key":    "secret_key",
							"session_token": "session_token",
							"region":        "x",
							"account_id":    "account_id",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Kind: "kind",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Kind: "kind",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type:          pkg.SourceTypeFileSystem,
						Configuration: map[string]string{},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
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
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
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
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
						DependsOn: []string{"source2"},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
						DependsOn: []string{"source2", "source3"},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
						Type: pkg.SourceTypeFileSystem,
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty relation criteria",
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{}},
			},
			wantErr: true,
		},
		{
			name: "missing relations criteria name",
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "missing relation criteria source.kind",
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "missing relation criteria source.meta_key",
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "missing relation criteria source.meta_value",
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:    "kind",
							MetaKey: "source_kind-key1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "missing relation criteria target.kind",
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "missing relation criteria target.meta_key",
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaValue: "related-metadata-value",
						},
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "missing relation criteria target.meta_value",
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
						Type: pkg.SourceTypeFileSystem,
						Configuration: map[string]string{
							"root_directory": "/root/path",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:    "related-kind",
							MetaKey: "related-metadata-key",
						},
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "missing GitHub token for GitHubRepository source",
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
						Type: pkg.SourceTypeGitHubRepository,
						Configuration: map[string]string{
							"user_or_org": "myuser",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing access key for AWS S3 source",
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
						Type: pkg.SourceTypeAWSS3,
						Configuration: map[string]string{
							"secret_key": "secret",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing secret key for AWS S3 source",
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
						Type: pkg.SourceTypeAWSS3,
						Configuration: map[string]string{
							"access_key":    "access",
							"session_token": "session_token",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing session token for AWS S3 source",
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
						Type: pkg.SourceTypeAWSS3,
						Configuration: map[string]string{
							"access_key": "access",
							"secret_key": "secret",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing region for AWS S3 source",
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
						Type: pkg.SourceTypeAWSS3,
						Configuration: map[string]string{
							"access_key":    "access",
							"secret_key":    "secret",
							"session_token": "session",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing access key for AWS S3 source",
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
						Type: pkg.SourceTypeAWSRDS,
						Configuration: map[string]string{
							"session_token": "session",
							"account_id":    "xxx",
							"region":        "xxxx",
							"secret_key":    "xxxx",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing account_id for AWS S3 source",
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
						Type: pkg.SourceTypeAWSRDS,
						Configuration: map[string]string{
							"session_token": "session",
							"region":        "xxxx",
							"secret_key":    "xxxx",
							"access_key":    "xxxx",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing secret key for AWS RDS source",
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
						Type: pkg.SourceTypeAWSRDS,
						Configuration: map[string]string{
							"access_key":    "access",
							"session_token": "session",
							"account_id":    "xxx",
							"region":        "xxxx",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing session token for AWS RDS source",
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
						Type: pkg.SourceTypeAWSRDS,
						Configuration: map[string]string{
							"access_key": "access",
							"secret_key": "secret",
							"account_id": "xxx",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing region for AWS RDS source",
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
						Type: pkg.SourceTypeAWSRDS,
						Configuration: map[string]string{
							"access_key":    "access",
							"secret_key":    "secret",
							"session_token": "session",
							"account_id":    "xxx",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing GitHub repository user_or_org",
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
						Type: pkg.SourceTypeGitHubRepository,
						Configuration: map[string]string{
							"token": "my-token",
						},
					},
				},
				Relation: Relation{RelationCriteria: []RelationCriteria{
					{
						Name: "name",
						Source: RelationCriteriaNode{
							Kind:      "kind",
							MetaKey:   "source_kind-key1",
							MetaValue: "source-kind-value1",
						},
						Target: RelationCriteriaNode{
							Kind:      "related-kind",
							MetaKey:   "related-metadata-key",
							MetaValue: "related-metadata-value",
						},
					},
				}},
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
