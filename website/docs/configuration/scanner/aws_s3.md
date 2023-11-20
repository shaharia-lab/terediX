---
title: "Configure AWS S3 Resource Scanner for TerediX"
sidebar_label: AWS S3
---

# AWS S3

<img src="/img/aws_s3_icon.png" alt="AWS S3" width="250"/>

## Configuration

### Type

Resource type. In this case it would be `aws_s3`.

### Configuration

- **access_key**: AWS access key
- **secret_key**: AWS secret key
- **region**: AWS region. e.g: us-west-1
- **session_token**: AWS session token
- **account_id**: AWS account ID

### Fields

List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data 
from the following fields.

- bucket_name
- region
- arn
- tags

### Schedule

**🔗 [Check schedule format](/docs/configuration/scanner/overview#schedule-format)**.

## Example

```yaml
source:
  aws_s3_source_one:
      type: aws_s3
      configuration:
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
      schedule: "@every 24h"
```

In the above example, we have added a source named `aws_s3_source_one` with type `aws_s3`. We have added some fields to add with each resource. 
We have also set the schedule to run this source every 24 hours.

Based on the above example, scanner_name would be `aws_s3_source_one` and scanner_type would be `aws_s3`. This is 
important to filter resources in Grafana dashboard.