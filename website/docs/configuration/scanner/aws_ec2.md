---
title: "AWS EC2"
---

# AWS EC2

<img src="/img/aws_ec2_icon.png" alt="AWS EC2" width="250"/>

## Configuration

### Type

Resource type. In this case it would be `aws_ec2`.

### Configuration

- **access_key**: AWS access key
- **secret_key**: AWS secret key
- **region**: AWS region. e.g: us-west-1
- **session_token**: AWS session token
- **account_id**: AWS account ID

### Fields

List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data 
from the following fields.

- instance_id
- image_id
- private_dns_name
- instance_type
- architecture
- instance_lifecycle
- instance_state
- vpc_id
- tags

### Schedule

**🔗 [Check schedule format](/docs/configuration/scanner/overview#schedule-format)**.

## Example

```yaml
source:
  aws_ec2_source_one:
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
      schedule: "@every 24h"
```

In the above example, we have added a source named `aws_ec2_source_one` with type `aws_ec2`. We have added some fields to add with each resource. 
We have also set the schedule to run this source every 24 hours.

Based on the above example, scanner_name would be `aws_ec2_source_one` and scanner_type would be `aws_ec2`. This is 
important to filter resources in Grafana dashboard.