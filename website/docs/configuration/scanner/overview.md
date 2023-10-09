---
sidebar_position: 1
title: "Overview"
---

# Source

Source is the place where terediX will discover the data. You can add multiple sources in the configuration file.

### Common source configuration

| option	                        | type           | description	                                                                                                              |
|--------------------------------|----------------|:--------------------------------------------------------------------------------------------------------------------------|
| 	  source_name                 | text           | Key/name of each source. e.g: github_repositories, aws_resources_rds, aws_s3_one	                                         |
| 	  [source_name].type          | text           | Type of the source. See the full list of supported source types	                                                          |
| 	  [source_name].configuration | key value pair | Configuration of the source. This configuration is different for different type of source	                                |
| 	  [source_name].fields        | list           | Additional data to store with each resource	                                                                              |
| 	  [source_name].schedule      | interval       | Set interval to schedule this source. e.g: `@every 10s`, `@every <br/>24h` or any valid cron expression `*/10 * * * * *`	 |

### Supported Source Types

| source type       | description	                                                                                                                                                 |
|-------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| aws_s3            | Discover data from AWS S3. You can use this source to discover data from AWS S3. See configuration for `aws_s3` source type.	                                |
| aws_ec2           | Discover data from AWS EC2. You can use this source to discover data from AWS EC2. See the configuration for `aws_ec2` source type.	                         |
| aws_rds           | Discover data from AWS RDS. You can use this source to discover data from AWS RDS. See the configuration for `aws_rds` source type.	                         |
| aws_ecr           | Discover data from AWS ECR. You can use this source to discover data from AWS ECR repository. See the configuration for `aws_ecr` source type.	              |
| file_system       | Discover data from local file system. You can use this source to discover data from local file system. See the configuration for `file_system` source type.	 |
| github_repository | List of GitHub repositories. See the configuration for `github_repository` source type                                                                       |
