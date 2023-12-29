package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewValidateCommand(t *testing.T) {
	// testCases holds the different test scenarios
	testCases := []struct {
		name           string         // name of the test case
		expectedOutput *cobra.Command // expected output of NewValidateCommand
		configFile     string         // input configFile
		expectedError  error          // error expected from validate() function
	}{
		{
			name: "valid_yml_file",
			expectedOutput: &cobra.Command{
				Use:   "validate",
				Short: "Validate YAML configuration",
				Long:  "Validate YAML configuration",
			},
			expectedError: nil,
			configFile:    "testdata/valid_config.yaml",
		},
		{
			name: "invalid_yml_file",
			expectedOutput: &cobra.Command{
				Use:   "validate",
				Short: "Validate YAML configuration",
				Long:  "Validate YAML configuration",
			},
			expectedError: errors.New("failed to unmarshal YAML data: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `this is...` into config.AppConfig"),
			configFile:    "testdata/invalid_config.yaml",
		},
		{
			name: "non_existing_file",
			expectedOutput: &cobra.Command{
				Use:   "validate",
				Short: "Validate YAML configuration",
				Long:  "Validate YAML configuration",
			},
			expectedError: errors.New("failed to read YAML file: open does_not_exist.yml: no such file or directory"),
			configFile:    "does_not_exist.yml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewValidateCommand()
			if cmd == nil {
				t.Fatal("Expected command, got nil")
				return
			}
			if cmd.Use != tc.expectedOutput.Use || cmd.Short != tc.expectedOutput.Short || cmd.Long != tc.expectedOutput.Long {
				t.Errorf("NewValidateCommand() => %v, want: %v", cmd, tc.expectedOutput)
			}
			err := validate(tc.configFile)
			if err != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	tests := []struct {
		name    string
		cfgFile string
		wantErr bool
	}{
		{
			name:    "valid config file",
			cfgFile: "testdata/valid_config.yaml",
			wantErr: false,
		},
		{
			name:    "invalid config file",
			cfgFile: "testdata/invalid_config.yaml",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.cfgFile); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
