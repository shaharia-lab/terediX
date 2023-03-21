package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
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

type AppConfig struct {
	Organization Organization      `yaml:"organization"`
	Discovery    Discovery         `yaml:"discovery"`
	Storage      Storage           `yaml:"storage"`
	Sources      map[string]Source `yaml:"source"`
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
	if c.Organization.Name == "" {
		return fmt.Errorf("organization name is required")
	}

	if c.Discovery.Name == "" {
		return fmt.Errorf("discovery name is required")
	}

	if c.Storage.BatchSize <= 0 {
		return fmt.Errorf("storage batch_size must be greater than 0")
	}

	if len(c.Storage.Engines) == 0 {
		return fmt.Errorf("at least one storage engine must be defined")
	}

	if _, ok := c.Storage.Engines[c.Storage.DefaultEngine]; !ok {
		return fmt.Errorf("default storage engine must be one of the defined engines")
	}

	// Validate storage engines
	for engine, config := range c.Storage.Engines {
		switch engine {
		case "postgresql":
			if err := c.validatePostgresqlEngine(config); err != nil {
				return err
			}
		case "neo4j":
			if err := c.validateNeo4jEngine(config); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown storage engine: '%s'", engine)
		}
	}

	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source must be defined")
	}

	for name, source := range c.Sources {
		if source.Type == "" {
			return fmt.Errorf("source '%s' type is required", name)
		}

		switch source.Type {
		case "file_system":
			if err := c.validateFileSystemSource(source); err != nil {
				return fmt.Errorf("source '%s': %v", name, err)
			}
		case "kubernetes":
			if err := c.validateKubernetesSource(source); err != nil {
				return fmt.Errorf("source '%s': %v", name, err)
			}
		default:
			return fmt.Errorf("unknown source type: '%s'", source.Type)
		}

		// Validate depends_on field
		for _, dependency := range source.DependsOn {
			if _, ok := c.Sources[dependency]; !ok {
				return fmt.Errorf("source '%s' depends_on contains invalid source key: '%s'", name, dependency)
			}
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

func (c *AppConfig) validateNeo4jEngine(config interface{}) error {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("neo4j engine configuration must be a map")
	}

	if _, ok := configMap["config_key"]; !ok {
		return fmt.Errorf("neo4j engine requires 'config_key'")
	}
	return nil
}

func (c *AppConfig) validateFileSystemSource(source Source) error {
	rootDirectory, ok := source.Configuration["root_directory"]
	if !ok || rootDirectory == "" {
		return fmt.Errorf("file_system source requires 'configuration.root_directory'")
	}
	return nil
}

func (c *AppConfig) validateKubernetesSource(source Source) error {
	kubeConfigFilePath, ok := source.Configuration["kube_config_file_path"]
	if !ok || kubeConfigFilePath == "" {
		return fmt.Errorf("kubernetes source requires 'configuration.kube_config_file_path'")
	}
	return nil
}
