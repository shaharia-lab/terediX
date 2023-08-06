package config

import (
	_ "embed"
	"errors"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v2"
)

//go:embed schema.json
var jsonSchema string

type SchemaValidator struct {
}

func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{}
}

func (sv *SchemaValidator) readJsonSchema() (string, error) {
	/*jsonSchema, err := jsonSchema.ReadFile("schema.json")
	if err != nil {
		return "", fmt.Errorf("failed to read JSON schema file: %w", err)
	}*/

	return jsonSchema, nil
}

func (sv *SchemaValidator) Validate(yamlContent string) (bool, error) {
	var m interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &m)
	if err != nil {
		panic(err)
	}
	m, err = sv.toStringKeys(m)
	if err != nil {
		panic(err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(jsonSchema)); err != nil {
		return false, err
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return false, err
	}

	if err := schema.Validate(m); err != nil {
		return false, err
	}

	return true, nil
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
