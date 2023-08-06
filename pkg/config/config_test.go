package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Test with valid config file
	path := "testdata/valid_config.yaml"
	_, err := Load(path)
	assert.NoError(t, err)

	// Test with non-existent config file
	path = "testdata/non_existent_config.yaml"
	_, err = Load(path)
	assert.Error(t, err)

	// Test with invalid YAML data
	path = "testdata/invalid_config.yaml"
	_, err = Load(path)
	assert.Error(t, err)
}
