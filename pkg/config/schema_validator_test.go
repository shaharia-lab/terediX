package config

import (
	"testing"
)

func TestSchemaValidator_Validate(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
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
  worker_pool_size: 1

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
    fields: &file_system_fields
      - rootDirectory
      - machineHost
    schedule: &schedule "@every 1m"
  github_repo:
    type: github_repository
    configuration:
      user_or_org: "some_org"
      token: "token"
    fields:
      - company
      - homepage
      - language
      - organization
      - stars
      - git_url
      - owner_login
      - owner_name
      - topics
    schedule: *schedule
  fs_two:
    type: file_system
    configuration:
      root_directory: "/some/other/path"
    fields: *file_system_fields
    schedule: *schedule
  aws_s3_one:
    type: aws_s3
    configuration: &aws_conf
      access_key: "xxxx"
      secret_key: "xxxx"
      session_token: "xxxx"
      region: "x"
      account_id: "xxx"
    fields:
      - bucket_name
      - region
      - arn
      - tags
    schedule: *schedule
  aws_rds_one:
    type: aws_rds
    configuration: *aws_conf
    fields:
      - instance_id
      - region
      - arn
      - tags
    schedule: *schedule
  aws_ec2_one:
    type: aws_ec2
    configuration:
      access_key: "xxxx"
      secret_key: "xxxx"
      session_token: "xxxx"
      region: "x"
      account_id: "xxx"
    fields:
      - instance_id
      - image_id
      - private_dns_name
      - instance_type
      - architecture
      - instance_lifecycle
      - instance_state
      - vpc_id
      - tags
    schedule: *schedule
  aws_ecr_example:
    type: aws_ecr
    configuration: *aws_conf
    fields:
      - repository_name
      - repository_uri
      - registry_id
      - arn
      - tags
    schedule: *schedule
relations:
  criteria:
    - name: "file-system-rule1"
      source:
        kind: "FilePath"
        meta_key: "rootDirectory"
        meta_value: "/some/path"
      target:
        kind: "FilePath"
        meta_key: "rootDirectory"
        meta_value: "/some/path"
`,
			expectError: false,
		},
		{
			name: "Invalid YAML against JSON schema",
			yamlContent: `
---
organization:
  name: My Org
  logo: http://example.com
discovery:
  description: ""
storage:
  batch_size: 2
  engines:
    postgresql:
      db: mydb
      host: localhost
      password: mypassword
      port: 5432
      user: myuser
  default_engine: postgresql
source:
  source1:
    type: file_system
    configuration:
      root_directory: /root/path
relations:
  criteria:
    - name: name
      kind: kind
      metadata_key: ""
      metadata_value: ""
      related_kind: ""
      related_metadata_key: ""
      related_metadata_value: ""
      source:
        kind: kind
        meta_key: source_kind-key1
        meta_value: source-kind-value1
      target:
        kind: related-kind
        meta_key: related-metadata-key
        meta_value: related-metadata-value
`,
			expectError: true,
		},
	}

	validator := NewSchemaValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.yamlContent)
			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if err != nil && !tt.expectError {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && tt.expectError {
				t.Errorf("expected error: %v", err)
			}
		})
	}
}
