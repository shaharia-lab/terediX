---
title: "Configure AWS ECR Resource Scanner for TerediX"
sidebar_label: AWS ECR
---

# AWS ECR

<img src="/img/aws_ecr_icon.png" alt="AWS ECR" width="250"/>

## Configuration

### Type

Resource type. In this case it would be `aws_ecr`.

### Configuration

- **access_key**: AWS access key
- **secret_key**: AWS secret key
- **region**: AWS region. e.g: us-west-1
- **session_token**: AWS session token
- **account_id**: AWS account ID

### Fields

List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data 
from the following fields.

- repository_name
- repository_uri
- registry_id
- arn
- tags

### Schedule

**🔗 [Check schedule format](/docs/configuration/scanner/overview#schedule-format)**.

## Example

```yaml
source:
  aws_ecr_source_one:
      type: aws_ecr
      configuration:
        access_key: "xxxx"
        secret_key: "xxxx"
        session_token: "xxxx"
        region: "x"
        account_id: "xxx"
      fields:
        - repository_name
        - repository_uri
        - registry_id
        - arn
        - tags
      schedule: "@every 24h"
```

In the above example, we have added a source named `aws_ecr_source_one` with type `aws_ecr`. We have added some fields to add with each resource. 
We have also set the schedule to run this source every 24 hours.

Based on the above example, scanner_name would be `aws_ecr_source_one` and scanner_type would be `aws_ecr`. This is 
important to filter resources in Grafana dashboard.