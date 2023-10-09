---
sidebar_position: 1
title: "Overview"
---

# Overview

**terediX** uses a configuration file to run. You can create a configuration file with the following command:

```bash
terediX init
```

It will generate the following configuration file.

```yaml
# config.yaml
---
organization:
  name: Your Organization
  logo: https://your-org-url.com/logo.png

discovery:
  name: Name of the discovery
  description: Some description about the discovery
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
  default_engine: postgresql

source:
  fs_one:
    type: file_system
    configuration:
      root_directory: "/root_directory"
    fields:
      - machineHost
      - rootDirectory
    schedule: &schedule "@every 1d"
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
```

This is the most basic one. You can add more sources, more storage engines, more discovery and more configuration as per your need.

## Validate Configuration

You can validate your configuration file with the following command:

```bash
teredix validate -c config.yaml
```

The above command will validate your `config.yaml` file against the [JSON schema](https://github.com/shaharia-lab/terediX/blob/master/pkg/config/schema.json).


## JSON Schema for Configuration

You can find the JSON schema for configuration file [here](https://github.com/shaharia-lab/terediX/blob/master/pkg/config/schema.json).