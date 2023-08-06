// Package cmd provides commands
package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/shahariaazam/teredix/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// NewValidateCommand build "discover" command
func NewValidateCommand() *cobra.Command {
	var cfgFile string

	cmd := cobra.Command{
		Use:   "validate",
		Short: "Validate YAML configuration",
		Long:  "Validate YAML configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return validate(cfgFile)
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "a valid yaml file is required")

	return &cmd
}

func validate(cfgFile string) error {
	yamlFile, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	var appConfig config.AppConfig
	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	schemaValidator := config.NewSchemaValidator()
	err = schemaValidator.Validate(string(yamlFile))
	if err != nil {
		return fmt.Errorf("invalid configuration. error: %w", err)
	}

	fmt.Println("Configuration is valid")

	return nil
}
