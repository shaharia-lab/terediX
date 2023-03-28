[![codecov](https://codecov.io/gh/shahariaazam/terediX/branch/master/graph/badge.svg?token=NKTKQ45HDN)](https://codecov.io/gh/shahariaazam/terediX)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=bugs)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX)

# teredix
`TeReDiX` (Tech Resource Discover &amp; Exploration) is a tool to discover tech resources for an organization and explore them as a resource graph

![Technical Resource Graph](https://user-images.githubusercontent.com/1095008/227896903-ff4b916d-893e-489e-8ffe-1469fff0703a.jpg)

## Technical Architecture

```
+-------------------+       +------------------------+        +------------------------+
| Source 1          |       | Scanner 1              |        | Storage Engine         |
| - Scanner         +------>+ - Scan()               |        | - Prepare()            |
|                   |       |                        |        | - Persist()            |
+-------------------+       +------------------------+        | - Find()               |
                                                                   +------------------------+
+-------------------+       +------------------------+        | - GetResources()       |
| Source 2          |       | Scanner 2              |        | - GetRelations()       |
| - Scanner         +------>+ - Scan()               |        | - StoreRelations()     |
|                   |       |                        |        +------------------------+
+-------------------+       +------------------------+

                           +------------------------+
                           | Processor              |
                           | - Process()            |
                           +------------------------+

                           +------------------------+
                           | Visualize              |
                           | - Render()             |
                           +------------------------+

```

The architecture consists of several components:

**Sources:** represent the different sources from which resources can be discovered. Each source has its own scanner to discover resources

**Scanners:** These are responsible for scanning a particular source and returning a list of resources.

**Storage Engine:** This component stores the discovered resources. It is responsible for preparing the storage schema, 
persisting resources, and finding resources based on a filter.

**Processor:** This component orchestrates the discovery process. It starts all the scanners in parallel and processes 
the resources as they become available.

**Visualizer:** This component is responsible for displaying the discovered resources and their relationships. It takes 
in a display type and renders the resources in that format.

## Installation

- Download the binary from latest release
- Run `chmod +x ./teredix` to make it executable

## Usage

```shell
➜ teredix --help
Discover and Explore your tech resources

Usage:
   [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  discover    Start discovering resources
  display     Display resource graph
  help        Help about any command
  relation    Build resource relationships

Flags:
  -h, --help      help for this command
  -v, --version   version for this command
```

## Config file

Here is an example `config.yaml` file. You should create your own config file with all the configuration.

```yaml
---
organization:
  name: Acme Tech
  logo: https://example.com

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
      db: "app
  default_engine: postgresql

source:
  fs_one:
    type: file_system
    configuration:
      root_directory: "/home/shaharia/Projects/teredix/pkg"
  fs_two:
    type: file_system
    configuration:
      root_directory: "/home/shaharia/Projects/k8s-resource-explorer"
  github_repo_sh:
    type: github_repository
    configuration:
      token: "xxxx"
      user_or_org: "xxxx"
```

| Field                                 | Type   | Required | Description                                                                                                                                                          |
|---------------------------------------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `sources`                             | map    | yes      | Configuration for the sources from where resources will be discovered.                                                                                               |
| `sources.<name>`                      | object | yes      | Configuration for a specific source.                                                                                                                                 |
| `sources.<name>.type`                 | string | yes      | The type of source. Currently, only `file_system` is supported.                                                                                                      |
| `sources.<name>.configuration`        | map    | yes      | The configuration for the source. The specific configuration options depend on the source type. See specific configuration for each source [here](#supported-source) |
| `storage`                             | object | yes      | Configuration for the storage where discovered resources will be saved.                                                                                              |
| `storage.engines`                     | map    | yes      | Configuration for the storage engines that can be used to save resources.                                                                                            |
| `storage.engines.postgresql`          | object | yes      | Configuration for the PostgreSQL engine.                                                                                                                             |
| `storage.engines.postgresql.host`     | string | yes      | The hostname of the PostgreSQL server.                                                                                                                               |
| `storage.engines.postgresql.port`     | int    | yes      | The port number of the PostgreSQL server.                                                                                                                            |
| `storage.engines.postgresql.user`     | string | yes      | The username to use when connecting to the PostgreSQL server.                                                                                                        |
| `storage.engines.postgresql.password` | string | yes      | The password to use when connecting to the PostgreSQL server.                                                                                                        |
| `storage.engines.postgresql.db`       | string | yes      | The name of the database to use on the PostgreSQL server.                                                                                                            |
| `processor`                           | object | yes      | Configuration for the processor that will discover resources and save them to storage.                                                                               |
| `processor.batch_size`                | int    | yes      | The number of resources to save to storage at once.                                                                                                                  |
| `log`                                 | object | no       | Configuration for logging.                                                                                                                                           |
| `log.level`                           | string | no       | The logging level to use. Valid values are `debug`, `info`, `warn`, and `error`. If not specified, the default logging level is `info`.                              |
| `log.format`                          | string | no       | The logging format to use. Valid values are `text` and `json`. If not specified, the default logging format is `text`.                                               |


## Supported Source

### GitHub Repository

It will fetch all GitHub repositories with metadata

```yaml
source:
  my_github_repositories:
    type: github_repository
    configuration:
      token: "xxx"
      user_or_org: "xxxx"
```

### File System

It will scan the file system and generate resource for each file. Required configuration:

```yaml
source:
  csv_data_sets:
    type: file_system
    configuration:
      root_directory: "/path/to/directory"
```

## Supported Storage

### PostgreSQL

To configure the PostgreSQL database, you need to update your configuration as following:

```yaml
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
```

`storage.batch_size` control how many resources should be inserted at once. Because all the scanner
run as goroutine and provide the resources as a channel for further processing/storing. So it's recommended
to use batch_size to avoid consuming heavy memory load if your organization has so many resources for all the source 
combined

## Supported Visualization

To visualize the resource graph, here is the available visualizer.

### Cytoscape JS

[Example](https://raw.githubusercontent.com/iVis-at-Bilkent/cytoscape.js-cise/master/demo.gif)

### Contributing

Contributions are welcome! Please follow the guidelines outlined in the [CONTRIBUTING](https://github.com/shahariaazam/teredix/blob/master/CONTRIBUTING.md) file.

### License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/shahariaazam/teredix/blob/master/LICENSE) file for details.