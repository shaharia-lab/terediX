// Package cmd provides a collection of commands for the Teredix application's command-line interface.
// These commands allow users to perform various tasks related to configuration validation.
//
// The commands are built using the spf13/cobra library, a powerful framework for building CLI applications in Go.
// Each command has its own logic defined in its corresponding function, and they utilize the pkg/config package
// for YAML configuration validation.
//
// The "validate" command validates a provided YAML configuration file against a predefined JSON schema.
// If the validation is successful, it prints a success message; otherwise, it reports validation errors.
//
// Example:
//   teredix validate -c config.yaml
//
// Dependencies:
//   - github.com/spf13/cobra: CLI command framework.
//   - github.com/shahariaazam/teredix/pkg/config: Package for YAML configuration validation.
//   - gopkg.in/yaml.v2: YAML parsing library.
//
// Usage:
//   func main() {
//     rootCmd := cmd.NewRootCommand()
//     if err := rootCmd.Execute(); err != nil {
//       fmt.Println(err)
//     }
//   }
//
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
