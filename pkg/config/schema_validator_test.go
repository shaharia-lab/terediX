package config

import (
	"testing"
)

func TestSchemaValidator_Validate(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		expected    bool
		expectError bool
	}{
		{
			name: "Valid YAML against JSON schema",
			yamlContent: `
---
organization:
  name: Shaharia Lab
  logo: http://example.com

discovery:
  name: Infrastructure Discovery
  description: Some description text

storage:
  batch_size: 2
  engines:
    postgresql:
      host: "localhost"
      port: 5432
      user: "app"
      password: "pass"
      db: "app"
    neo4j:
      config_key: "value"
  default_engine: postgresql

source:
  fs_one:
    type: file_system
    configuration:
      root_directory: "/some/path"
  fs_two:
    type: file_system
    configuration:
      root_directory: "/some/other/path"
  aws_s3_one:
    type: aws_s3
    configuration:
      access_key: "xxxx"
      secret_key: "xxxx"
      session_token: "xxxx"
      region: "x"
      account_id: "xxx"
  aws_rds_one:
    type: aws_rds
    config_from: aws_s3_one
  aws_ec2_one:
    type: aws_ec2
    config_from: aws_s3_one
  aws_ecr_example:
    type: aws_ecr
    config_from: aws_s3_one
relations:
  criteria:
    - name: "file-system-rule1"
      source:
        kind: "FilePath"
        meta_key: "Root-Directory"
        meta_value: "/some/path"
      target:
        kind: "FilePath"
        meta_key: "Root-Directory"
        meta_value: "/some/path"
`,
			expected:    true,
			expectError: false,
		},
	}

	validator := NewSchemaValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := validator.Validate(tt.yamlContent)
			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if valid != tt.expected {
				t.Errorf("expected valid: %v, but got: %v", tt.expected, valid)
			}
		})
	}
}
