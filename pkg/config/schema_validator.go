/*
Package config provides functionality to validate YAML configuration using a JSON schema.

This package defines a SchemaValidator type that allows you to validate a YAML content against a predefined JSON schema.
The JSON schema is embedded using the //go:embed directive.

Usage:
	validator := config.NewSchemaValidator()

	yamlContent := `
		source:
		  fs_one:
		    type: file_system
		    configuration:
		      root_directory: "/some/path"
		...`

	err := validator.Validate(yamlContent)
	if err != nil {
		fmt.Println("Validation error:", err)
	}

Note:
	The embedded JSON schema is expected to be stored in a file named "schema.json" in the same package directory.

For more details about the embedded JSON schema, refer to the embedded schema documentation.

This package relies on the santhosh-tekuri/jsonschema and gopkg.in/yaml.v2 libraries.
*/
package config

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v2"
)

//go:embed schema.json
var jsonSchema string

// SchemaValidator struct represent schema validator
type SchemaValidator struct {
}

// NewSchemaValidator construct a new validator
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{}
}

// Validate validates the provided YAML content against the embedded JSON schema.
// It checks whether the provided YAML conforms to the expected structure and data types defined in the JSON schema.
// If the validation fails, an error is returned indicating the reason for the failure.
//
// Parameters:
//   - yamlContent: A string containing the YAML content to be validated.
//
// Returns:
//   - An error indicating the validation failure, if any.
//
// Example:
//   validator := config.NewSchemaValidator()
//   yamlContent := `
//     source:
//       fs_one:
//         type: file_system
//         configuration:
//           root_directory: "/some/path"
//     ...`
//   err := validator.Validate(yamlContent)
//   if err != nil {
//     fmt.Println("Validation error:", err)
//   }
//
// Note:
//   - The embedded JSON schema is expected to be stored in a file named "schema.json" in the same package directory.
//   - This method relies on the santhosh-tekuri/jsonschema and gopkg.in/yaml.v2 libraries.
//
func (sv *SchemaValidator) Validate(yamlContent string) error {
	var m interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &m)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml content. Error: %w", err)
	}
	m, err = sv.toStringKeys(m)
	if err != nil {
		return fmt.Errorf("failed to prepare yaml content for validation against json schema. Error: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(jsonSchema)); err != nil {
		return fmt.Errorf("json schema compiler failed. Error: %w", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("failed to compile json schema during validation. Error: %w", err)
	}

	if err := schema.Validate(m); err != nil {
		return fmt.Errorf("failed to validate configuration. Error: %w", err)
	}

	return nil
}

func (sv *SchemaValidator) toStringKeys(val interface{}) (interface{}, error) {
	var err error
	switch val := val.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			k, ok := k.(string)
			if !ok {
				return nil, errors.New("found non-string key")
			}
			m[k], err = sv.toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case []interface{}:
		var l = make([]interface{}, len(val))
		for i, v := range val {
			l[i], err = sv.toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return val, nil
	}
}
