// Package config read and validate config file
package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Organization store organization data
type Organization struct {
	Name string `yaml:"name"`
	Logo string `yaml:"logo"`
}

// Discovery hold discovery configuration
type Discovery struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Storage store storage configuration
type Storage struct {
	BatchSize     int                    `yaml:"batch_size"`
	Engines       map[string]interface{} `yaml:"engines"`
	DefaultEngine string                 `yaml:"default_engine"`
}

// SourceConfiguration store source configuration data
type SourceConfiguration struct {
	RootDirectory string `yaml:"root_directory"`
}

// Source holds source configuration
type Source struct {
	Type          string            `yaml:"type"`
	Configuration map[string]string `yaml:"configuration"`
	Fields        []string          `yaml:"fields"`
	Schedule      string            `yaml:"schedule,omitempty"`
}

// RelationCriteriaNode represent source and target
type RelationCriteriaNode struct {
	Kind      string `yaml:"kind"`
	MetaKey   string `yaml:"meta_key"`
	MetaValue string `yaml:"meta_value"`
}

// RelationCriteria represents criteria for relation builder
type RelationCriteria struct {
	Name                 string               `yaml:"name"`
	Kind                 string               `yaml:"kind"`
	MetadataKey          string               `yaml:"metadata_key"`
	MetadataValue        string               `yaml:"metadata_value"`
	RelatedKind          string               `yaml:"related_kind"`
	RelatedMetadataKey   string               `yaml:"related_metadata_key"`
	RelatedMetadataValue string               `yaml:"related_metadata_value"`
	Source               RelationCriteriaNode `yaml:"source"`
	Target               RelationCriteriaNode `yaml:"target"`
}

// Relation represent relationship rules
type Relation struct {
	RelationCriteria []RelationCriteria `yaml:"criteria"`
}

// AppConfig provides configuration for the tools
type AppConfig struct {
	Organization Organization      `yaml:"organization"`
	Discovery    Discovery         `yaml:"discovery"`
	Storage      Storage           `yaml:"storage"`
	Sources      map[string]Source `yaml:"source"`
	Relation     Relation          `yaml:"relations"`
}

// Load loads configuration file
func Load(path string) (*AppConfig, error) {
	var appConfig AppConfig

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return &appConfig, fmt.Errorf("failed to read YAML file: %w", err)
	}

	schemaValidator := NewSchemaValidator()
	err = schemaValidator.Validate(string(yamlFile))
	if err != nil {
		return nil, fmt.Errorf("invalid configuration file %s. Error: %w", path, err)
	}

	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		return &appConfig, fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	sourceConfigs := map[string]map[string]string{}
	for sourceName, s := range appConfig.Sources {
		sourceConfigs[sourceName] = s.Configuration
	}

	return &appConfig, nil
}
