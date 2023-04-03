<h1 align="center">TerediX</h1>
<p align="center">Tech Resource Graph (Discovery & Exploration)</p>

<p align="center">
  <a href="https://github.com/shahariaazam/terediX/actions/workflows/CI.yaml"><img src="https://github.com/shahariaazam/terediX/actions/workflows/CI.yaml/badge.svg" height="20"/></a>
  <a href="https://codecov.io/gh/shahariaazam/terediX"><img src="https://codecov.io/gh/shahariaazam/terediX/branch/master/graph/badge.svg?token=NKTKQ45HDN" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=reliability_rating" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=vulnerabilities" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=security_rating" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=sqale_rating" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=code_smells" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=ncloc" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=alert_status" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=duplicated_lines_density" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=bugs" height="20"/></a>
  <a href="https://sonarcloud.io/summary/new_code?id=shahariaazam_terediX"><img src="https://sonarcloud.io/api/project_badges/measure?project=shahariaazam_terediX&metric=sqale_index" height="20"/></a>
</p><br/><br/>

<p align="center">
  <a href="https://github.com/shahariaazam/teredix"><img src="https://user-images.githubusercontent.com/1095008/229536376-51ddaa75-85ee-4e3e-95df-7cf6093392bf.png" width="100%"/></a>
</p><br/>


## 🤔  What is TerediX?

`TeReDiX` (Tech Resource Discover &amp; Exploration) is a tool to discover tech resources for an organization and
explore them as a resource graph. Users can use **teredix** to create resource graphs that depict the relationships
and dependencies between different resources. The tool also allows users to search for specific resources and view details about each resource, such as its name, type, description, and associated tags.

**teredix** can be a useful tool for organizations looking to manage their technology resources more effectively and
gain a better understanding of how those resources are interconnected.

This tools can efficiently collect all the resources and their relevent metadata from different sources and can build
relationship between resources to enhance visibility of resources.

## 📄 Table of Contents
- [Technical Architecture](#-technical-architecture)
- [Getting Started](#-getting-started)
  - [Run Discovery](#-run-discovery)
  - [Build the Resource Relationship](#-build-the-resource-relationship)
  - [Explore Resource Graph Visualization](#-explore-resource-graph-visualization)
- [Usage](#-usage)
- [Config File](#-config-file)
- [Supported Source](#-supported-source)
  - [GitHub Repository](#octocat-github-repository)
  - [File System](#file-system)
- [Supported Storage](#supported-storage)
  - [PostgreSQL](#postgresql)
- [Supported Visualization](#supported-visualization)
  - [Cytoscape JS](#cytoscape-js)
- [Contributing](#contributing)
- [License](#license)


## 💻  Technical Architecture

```puml
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

**Storage Engine:** This component stores the discovered resources. It is responsible for preparing the storage
schema, persisting resources, and finding resources based on a filter.

**Processor:** This component orchestrates the discovery process. It starts all the scanners in parallel and
processes the resources as they become available.

**Visualizer:** This component is responsible for displaying the discovered resources and their relationships. It
takes in a display type and renders the resources in that format.

## 🚀 Getting Started

- To get started, you just need to download the latest binary from [latest release](https://github.com/shahariaazam/terediX/releases).
- Prepare your config file. Here is the [example config file](#-config-file).

**This tool has three major commands available**

| Command                  | Description                                                                         |
|--------------------------|-------------------------------------------------------------------------------------|
| discover                 | It will go through each sources and discover resources and save it to your database |
| relation                 | This command is responsible to build relationship between resources                 |
| display                  | It will generate visualized resource graph that can be opened in your browser       |

### 📥 Run Discovery

To discover the resources, run the following command.

```shell
teredix discover --config {your_config.yaml file}
```

### 🔗 Build the resource relationship

If you want to build the relationship between resources based on your relation criteria defined in your config file,
run -

```shell
teredix relation --config {your_config.yaml file}
```

### 🔍 Explore Resource Graph Visualization

To display the resource graph in an interactive way, please run the following command.

```shell
teredix display --config {your_config.yaml file}
```
It will show `Displaying resource graph at http://localhost:8989`. Open your browser and visit the URL.

## 💻 Usage

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

## ⚙️ Config file

Here is an example `config.yaml` file. You should create your own config file with all the configuration.

```yaml
---
organization:
  name: "Acme Tech"
  logo: "https://example.com"
discovery:
  name: "Infrastructure Discovery"
  description: "Some description text"
storage:
  batch_size: 2
  engines:
    postgresql:
      host: "localhost"
      port: 5432
      user: "app"
      password: "pass"
      db: "app"
  default_engine: "postgresql"
source:
  fs_one:
    type: "file_system"
    configuration:
      root_directory: "/some/path"
  fs_two:
    type: "file_system"
    configuration:
      root_directory: "/some/other/path"
  github_repo_sh:
    type: "github_repository"
    configuration:
      token: "xxxx"
      user_or_org: "xxxx"
  aws_s3_one:
    type: "aws_s3"
    configuration:
      access_key: "xxxx"
      secret_key: "xxxx"
      session_token: "xxxx"
      region: "x"
      account_id: "xxx"
  aws_rds_one:
    type: "aws_rds"
    config_from: "aws_s3_one"
  aws_ec2_org:
    type: "aws_ec2"
    config_from: "aws_s3_one"
  aws_ecr_org:
    type: "aws_ecr"
    config_from: "aws_s3_one"
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
        meta_value: "/some/other/path"
```

You can get a complete valid config file in [valid_config.yaml](https://github.com/shahariaazam/terediX/blob/master/pkg/config/testdata/valid_config.yaml)


| Field                                 | Type   | Required | Description                                                                                                                                                           |
|---------------------------------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `sources`                             | map    | yes      | Configuration for the sources from where resources will be discovered.                                                                                                |
| `sources.<name>`                      | object | yes      | Configuration for a specific source.                                                                                                                                  |
| `sources.<name>.type`                 | string | yes      | The type of source. Currently, only `file_system` is supported.                                                                                                       |
| `sources.<name>.configuration`        | map    | yes      | The configuration for the source. The specific configuration options depend on the source type. See specific configuration for each source [here](#-supported-source) |
| `storage`                             | object | yes      | Configuration for the storage where discovered resources will be saved.                                                                                               |
| `storage.engines`                     | map    | yes      | Configuration for the storage engines that can be used to save resources.                                                                                             |
| `storage.engines.postgresql`          | object | yes      | Configuration for the PostgreSQL engine.                                                                                                                              |
| `storage.engines.postgresql.host`     | string | yes      | The hostname of the PostgreSQL server.                                                                                                                                |
| `storage.engines.postgresql.port`     | int    | yes      | The port number of the PostgreSQL server.                                                                                                                             |
| `storage.engines.postgresql.user`     | string | yes      | The username to use when connecting to the PostgreSQL server.                                                                                                         |
| `storage.engines.postgresql.password` | string | yes      | The password to use when connecting to the PostgreSQL server.                                                                                                         |
| `storage.engines.postgresql.db`       | string | yes      | The name of the database to use on the PostgreSQL server.                                                                                                             |
| `processor`                           | object | yes      | Configuration for the processor that will discover resources and save them to storage.                                                                                |
| `processor.batch_size`                | int    | yes      | The number of resources to save to storage at once.                                                                                                                   |
| `log`                                 | object | no       | Configuration for logging.                                                                                                                                            |
| `log.level`                           | string | no       | The logging level to use. Valid values are `debug`, `info`, `warn`, and `error`. If not specified, the default logging level is `info`.                               |
| `log.format`                          | string | no       | The logging format to use. Valid values are `text` and `json`. If not specified, the default logging format is `text`.                                                |


The following fields can be specified for each criteria:

| Field               | Type   | Required | Description                                       |
|---------------------|--------|----------|---------------------------------------------------|
| `name`              | string | yes      | A descriptive name for the criteria.              |
| `source`            | object | yes      | Configuration for the source of the relationship. |
| `source.kind`       | string | yes      | The kind of resource to analyze.                  |
| `source.meta_key`   | string | yes      | The metadata key used to identify the resource.   |
| `source.meta_value` | string | yes      | The metadata value used to identify the resource. |
| `target`            | object | yes      | Configuration for the target of the relationship. |
| `target.kind`       | string | yes      | The kind of resource to analyze.                  |
| `target.meta_key`   | string | yes      | The metadata key used to identify the resource.   |
| `target.meta_value` | string | yes      | The metadata value used to identify the resource. |

## 🌐 Supported Source

### AWS EC2

Discover all EC2 instances and it's tags

```yaml
source:
  aws_ec2_first_config:
    type: aws_ec2
    configuration:
      access_key: "xxxx"
      secret_key: "xxxx"
      session_token: "xxxx"
      region: "x"
      account_id: "xxx
```

#### Available metadata for AWS EC2 Instance

| Meta Key                                 | Description                                                       |
|------------------------------------------|-------------------------------------------------------------------|
| AWS-EC2-Instance-ID                      | The ID of the instance.                                           |
| AWS-EC2-Image-ID                         | The ID of the AMI                                                 |
| AWS-EC2-PrivateDnsName                   | The private IPv4 DNS name                                         |
| AWS-EC2-InstanceType                     | The instance type                                                 |
| AWS-EC2-Architecture                     | The architecture of the instance                                  |
| AWS-EC2-InstanceLifecycle                | Indicates whether this is a Spot Instance or a Scheduled Instance |
| AWS-EC2-InstanceState                    | The current state of the instance                                 |
| AWS-EC2-VpcId                            | The ID of the VPC in which the instance is running                |
| AWS-EC2-Tag-{tagKey}                     | Any tags assigned to the instance                                 |
| AWS-EC2-Security-Group-{securityGroupID} | Security group ID                                                 |
| Scanner-Label                            | Name of the source configured in config.yaml file                 |


### AWS ECR

Discover all ECR repositories and it's tags

```yaml
source:
  aws_ecr_work_config:
    type: aws_ecr
    configuration:
      access_key: "xxxx"
      secret_key: "xxxx"
      session_token: "xxxx"
      region: "x"
      account_id: "xxx
```

#### Available metadata for AWS ECR repository

| Meta Key                | Description                                                                |
|-------------------------|----------------------------------------------------------------------------|
| AWS-ECR-Registry-Id     | ECR registry ID.                                                           |
| AWS-ECR-Repository-Name | ECR Repository name. eg. hello-world                                       |
| AWS-ECR-Repository-URI  | Repository URI. eg. 123456789.dkr.ecr.eu-west-1.amazonaws.com/hello-world  |
| AWS-ECR-Repository-Arn  | Repository ARN. eg. arn:aws:ecr:eu-west-1:123456789:repository/hello-world |
| AWS-ECR-{tagKey}        | Tags of ECR repositories                                                   |
| Scanner-Label           | Name of the source configured in config.yaml file                          |


### AWS RDS

Discover all RDS instances and it's tags

```yaml
source:
  aws_rds_one:
    type: aws_rds
    configuration:
      access_key: "xxxx"
      secret_key: "xxxx"
      session_token: "xxxx"
      region: "x"
      account_id: "xxx
```

#### Available metadata for AWS RDS Instance

| Meta Key             | Description                                                                                          |
|----------------------|------------------------------------------------------------------------------------------------------|
| AWS-RDS-Instance-ID  | Name of the RDS instance                                                                             |
| Scanner-Label        | Name of the source configured in config.yaml file                                                    |
| AWS-RDS-Region       | Region name of AWS RDS source                                                                        |
| AWS-ARN              | AWS ARN for the RDS instance                                                                         |
| AWS-RDS-Tag-{tagKey} | Every tag key associated with the source will be added as metadata. {tagKey} is the tag key from RDS |


### AWS S3

Discover all the S3 buckets from AWS

```yaml
source:
  my_org_s3_bucket_eu:
    type: aws_s3
    configuration:
      access_key: "xxxx"
      secret_key: "xxxx"
      session_token: "xxxx"
      region: "x"
      account_id: "xxx"
```

#### Available metadata for AWS S3 Source

| Meta Key            | Description                                                                                         |
|---------------------|-----------------------------------------------------------------------------------------------------|
| AWS-S3-Bucket-Name  | Name of the bucket name                                                                             |
| Scanner-Label       | Name of the source configured in config.yaml file                                                   |
| AWS-S3-Region       | Region name of AWS S3 source                                                                        |
| AWS-ARN             | AWS ARN for the bucket                                                                              |
| AWS-S3-Tag-{tagKey} | Every tag key associated with the source will be added as metadata. {tagKey} is the tag key from S3 |


### :octocat: GitHub Repository

It will fetch all GitHub repositories with metadata

```yaml
source:
  my_github_repositories:
    type: github_repository
    configuration:
      token: "xxx"
      user_or_org: "xxxx"
```

#### Available metadata for github_repository source

| Meta Key                 | Description                                       |
|--------------------------|---------------------------------------------------|
| GitHub-Repo-Language     | Primary language detected for this repository     |
| GitHub-Repo-Stars        | Total count of stars                              |
| Scanner-Label            | Name of the source configured in config.yaml file |
| GitHub-Repo-Git-URL      | Git URL for the repository                        |
| GitHub-Repo-Homepage     | Homepage URL of GitHub repository                 |
| GitHub-Repo-Organization | Name of the organization                          |
| GitHub-Owner             | Owner of GitHub repository                        |
| GitHub-Company           | Company name of the GitHub repository owner       |
| GitHub-Repo-Topic        | Topics of GitHub repository                       |


### File System

It will scan the file system and generate resource for each file. Required configuration:

```yaml
source:
  csv_data_sets:
    type: file_system
    configuration:
      root_directory: "/path/to/directory"
```
#### Available metadata for file_system source

| Meta Key       | Description                                       |
|----------------|---------------------------------------------------|
| Machine-Host   | Hostname of the machine                           |
| Root-Directory | Root directory to scan the recursive file list    |
| Scanner-Label  | Name of the source configured in config.yaml file |


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