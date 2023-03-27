package config

import (
	"fmt"
	"io/ioutil"
	"teredix/pkg"

	"gopkg.in/yaml.v3"
)

type Organization struct {
	Name string `yaml:"name"`
	Logo string `yaml:"logo"`
}

type Discovery struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Storage struct {
	BatchSize     int                    `yaml:"batch_size"`
	Engines       map[string]interface{} `yaml:"engines"`
	DefaultEngine string                 `yaml:"default_engine"`
}

type SourceConfiguration struct {
	RootDirectory string `yaml:"root_directory"`
}

type Source struct {
	Type          string            `yaml:"type"`
	Configuration map[string]string `yaml:"configuration"`
	DependsOn     []string          `yaml:"depends_on,omitempty"`
}

type RelationCriteria struct {
	Name                 string `yaml:"name"`
	Kind                 string `yaml:"kind"`
	MetadataKey          string `yaml:"metadata_key"`
	MetadataValue        string `yaml:"metadata_value"`
	RelatedKind          string `yaml:"related_kind"`
	RelatedMetadataKey   string `yaml:"related_metadata_key"`
	RelatedMetadataValue string `yaml:"related_metadata_value"`
}

type Relation struct {
	RelationCriteria []RelationCriteria `yaml:"criteria"`
}

type AppConfig struct {
	Organization Organization      `yaml:"organization"`
	Discovery    Discovery         `yaml:"discovery"`
	Storage      Storage           `yaml:"storage"`
	Sources      map[string]Source `yaml:"source"`
	Relation     Relation          `yaml:"relations"`
}

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

	err = c.validateStorage(c.Storage, c.Sources)
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

func (c *AppConfig) validateStorage(storage Storage, sources map[string]Source) error {
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

	if len(sources) == 0 {
		return fmt.Errorf("at least one source must be defined")
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
		if err := c.validateFileSystemSourceConfiguration(name, source); err != nil {
			return err
		}
	case "kubernetes":
		if err := c.validateKubernetesSourceConfiguration(name, source); err != nil {
			return err
		}
	case pkg.SourceTypeGitHubRepository:
		if err := c.validateGitHubRepositorySourceConfiguration(name, source); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown source type: '%s'", source.Type)
	}
	return nil
}

func (c *AppConfig) validateFileSystemSourceConfiguration(name string, source Source) error {
	rootDirectory, ok := source.Configuration["root_directory"]
	if !ok || rootDirectory == "" {
		return fmt.Errorf("source '%s' requires 'configuration.root_directory'", name)
	}
	return nil
}

func (c *AppConfig) validateKubernetesSourceConfiguration(name string, source Source) error {
	kubeConfigFilePath, ok := source.Configuration["kube_config_file_path"]
	if !ok || kubeConfigFilePath == "" {
		return fmt.Errorf("source '%s' requires 'configuration.kube_config_file_path'", name)
	}
	return nil
}

func (c *AppConfig) validateGitHubRepositorySourceConfiguration(name string, source Source) error {
	ghToken, ok := source.Configuration["token"]
	if !ok || ghToken == "" {
		return fmt.Errorf("source '%s' requires 'configuration.token'", name)
	}

	userOrOrg, ok := source.Configuration["user_or_org"]
	if !ok || userOrOrg == "" {
		return fmt.Errorf("source '%s' requires 'configuration.user_or_org'", name)
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
	}

	return nil
}
