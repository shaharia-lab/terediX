// Package config read and validate config file
package config

import (
	"fmt"
	"io/ioutil"
	"teredix/pkg"

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
	DependsOn     []string          `yaml:"depends_on,omitempty"`
}

// RelationCriteria represents criteria for relation builder
type RelationCriteria struct {
	Name                 string `yaml:"name"`
	Kind                 string `yaml:"kind"`
	MetadataKey          string `yaml:"metadata_key"`
	MetadataValue        string `yaml:"metadata_value"`
	RelatedKind          string `yaml:"related_kind"`
	RelatedMetadataKey   string `yaml:"related_metadata_key"`
	RelatedMetadataValue string `yaml:"related_metadata_value"`
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

	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		return &appConfig, fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	return &appConfig, nil
}

// Validate Add this method to your AppConfig struct
func Validate(c *AppConfig) error {
	err := c.validateOrganization(c.Organization)
	if err != nil {
		return err
	}

	err = c.validateDiscovery(c.Discovery)
	if err != nil {
		return err
	}

	err = c.validateStorage(c.Storage)
	if err != nil {
		return err
	}

	err = c.validateSources(c.Sources)
	if err != nil {
		return err
	}

	err = c.validateRelations(c.Relation)
	if err != nil {
		return err
	}

	return nil
}

func (c *AppConfig) validateOrganization(org Organization) error {
	if org.Name == "" {
		return fmt.Errorf("organization name is required")
	}

	return nil
}

func (c *AppConfig) validateDiscovery(discovery Discovery) error {
	if discovery.Name == "" {
		return fmt.Errorf("discovery name is required")
	}

	return nil
}

func (c *AppConfig) validateStorage(storage Storage) error {
	if storage.BatchSize <= 0 {
		return fmt.Errorf("storage batch_size must be greater than 0")
	}

	if len(storage.Engines) == 0 {
		return fmt.Errorf("at least one storage engine must be defined")
	}

	if _, ok := storage.Engines[storage.DefaultEngine]; !ok {
		return fmt.Errorf("default storage engine must be one of the defined engines")
	}

	// Validate storage engines
	for engine, config := range storage.Engines {
		switch engine {
		case "postgresql":
			if err := c.validatePostgresqlEngine(config); err != nil {
				return err
			}
		case "neo4j":
			if err := validateNeo4jEngine(config); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown storage engine: '%s'", engine)
		}
	}

	return nil
}

func (c *AppConfig) validatePostgresqlEngine(config interface{}) error {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("postgresql engine configuration must be a map")
	}

	requiredKeys := []string{"host", "user", "password", "db", "port"}
	for _, key := range requiredKeys {
		if _, ok := configMap[key]; !ok {
			return fmt.Errorf("postgresql engine requires '%s'", key)
		}
	}
	return nil
}

func validateNeo4jEngine(config interface{}) error {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("neo4j engine configuration must be a map")
	}

	if _, ok := configMap["config_key"]; !ok {
		return fmt.Errorf("neo4j engine requires 'config_key'")
	}
	return nil
}

func (c *AppConfig) validateSources(sources map[string]Source) error {
	if len(sources) == 0 {
		return fmt.Errorf("at least one source must be defined")
	}

	for name, source := range sources {
		if source.Type == "" {
			return fmt.Errorf("source '%s' type is required", name)
		}

		if err := c.validateSourceConfiguration(source.Type, source); err != nil {
			return fmt.Errorf("source '%s': %v", name, err)
		}

		if err := c.validateDependsOn(name, source); err != nil {
			return fmt.Errorf("source '%s': %v", name, err)
		}
	}

	return nil
}

func (c *AppConfig) validateSourceConfiguration(name string, source Source) error {
	switch source.Type {
	case pkg.SourceTypeFileSystem:
		if err := c.validateConfigurationKeys(name, source, "root_directory"); err != nil {
			return err
		}
	case "kubernetes":
		if err := c.validateConfigurationKeys(name, source, "kube_config_file_path"); err != nil {
			return err
		}
	case pkg.SourceTypeGitHubRepository:
		if err := c.validateConfigurationKeys(name, source, "token", "user_or_org"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown source type: '%s'", source.Type)
	}
	return nil
}

func (c *AppConfig) validateConfigurationKeys(sourceName string, source Source, requiredKeys ...string) error {
	for _, k := range requiredKeys {
		keyNotEmpty, ok := source.Configuration[k]
		if !ok || keyNotEmpty == "" {
			return fmt.Errorf("source '%s' requires 'configuration.%s'", k, sourceName)
		}
	}

	return nil
}

func (c *AppConfig) validateDependsOn(name string, source Source) error {
	for _, dependency := range source.DependsOn {
		if _, ok := c.Sources[dependency]; !ok {
			return fmt.Errorf("source '%s' depends_on contains invalid source key: '%s'", name, dependency)
		}
	}
	return nil
}

func (c *AppConfig) validateRelations(relations Relation) error {
	if relations.RelationCriteria == nil {
		return fmt.Errorf("relations field must be defined")
	}

	if len(relations.RelationCriteria) == 0 {
		return fmt.Errorf("relations.criteria is empty")
	}

	for _, criteria := range relations.RelationCriteria {
		err := c.validateRelationCriteria(criteria)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *AppConfig) validateRelationCriteria(criteria RelationCriteria) error {
	if criteria.Name == "" {
		return fmt.Errorf("relations.criteria.name is required")
	}

	if criteria.Kind == "" {
		return fmt.Errorf("relations.criteria.kind is required")
	}

	if criteria.MetadataKey == "" {
		return fmt.Errorf("relations.criteria.metadata_key is required")
	}

	if criteria.MetadataValue == "" {
		return fmt.Errorf("relations.criteria.metadata_value is required")
	}

	if criteria.RelatedKind == "" {
		return fmt.Errorf("relations.criteria.related_kind is required")
	}

	if criteria.RelatedMetadataKey == "" {
		return fmt.Errorf("relations.criteria.related_metadata_key is required")
	}

	if criteria.RelatedMetadataValue == "" {
		return fmt.Errorf("relations.criteria.related_metadata_value is required")
	}
	return nil
}
