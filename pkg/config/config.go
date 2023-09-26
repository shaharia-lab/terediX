// Package config read and validate config file
package config

import (
	"fmt"
	"io/ioutil"

	"github.com/shaharia-lab/teredix/pkg"

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
	ConfigFrom    string            `yaml:"config_from,omitempty"`
	Configuration map[string]string `yaml:"configuration"`
	Fields        []string          `yaml:"fields"`
	DependsOn     []string          `yaml:"depends_on,omitempty"`
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

	for name, source := range appConfig.Sources {
		if source.ConfigFrom != "" && sourceConfigs[source.ConfigFrom] != nil {
			sourceConfiguration := sourceConfigs[source.ConfigFrom]
			source.Configuration = sourceConfiguration
			appConfig.Sources[name] = source
		}
	}

	return &appConfig, nil
}

func (c *AppConfig) validateDiscovery(discovery Discovery) error {
	if discovery.Name == "" {
		return fmt.Errorf("discovery name is required")
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
	case pkg.SourceTypeAWSS3, pkg.SourceTypeAWSRDS, pkg.SourceTypeAWSEC2, pkg.SourceTypeAWSECR:
		if err := c.validateConfigurationKeys(name, source, "access_key", "secret_key", "session_token", "region", "account_id"); err != nil {
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
			return fmt.Errorf("source '%s' requires 'configuration.%s'", sourceName, k)
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

	if criteria.Source.Kind == "" {
		return fmt.Errorf("relations.criteria.source.kind is required")
	}

	if criteria.Source.MetaKey == "" {
		return fmt.Errorf("relations.criteria.source.meta_key is required")
	}

	if criteria.Source.MetaValue == "" {
		return fmt.Errorf("relations.criteria.source.meta_value is required")
	}

	if criteria.Target.Kind == "" {
		return fmt.Errorf("relations.criteria.target.kind is required")
	}

	if criteria.Target.MetaKey == "" {
		return fmt.Errorf("relations.criteria.target.meta_key is required")
	}

	if criteria.Target.MetaValue == "" {
		return fmt.Errorf("relations.criteria.target.meta_value is required")
	}
	return nil
}
