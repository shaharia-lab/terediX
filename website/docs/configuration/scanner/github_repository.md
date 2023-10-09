---
title: "GitHub Repository"
---

# GitHub Repository

<img src="/img/github_repository_icon.png" alt="GitHub Repository" width="250"/>

## Configuration

### Type

Resource type. In this case it would be `github_repository`.

### Configuration

- **user_or_org**: User or organization username of GitHub.
- **token**: GitHub access token

### Fields

List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data 
from the following fields.

- company
- homepage
- language
- organization
- stars
- git_url
- owner_login
- owner_name
- topics

### Schedule

**🔗 [Check schedule format](/docs/configuration/scanner/overview#schedule-format)**.

## Example

```yaml
source:
  github_repository_source_one:
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
      schedule: "@every 24h"
```

In the above example, we have added a source named `github_repository_source_one` with type `github_repository`. We have added some fields to add with each resource. 
We have also set the schedule to run this source every 24 hours.

Based on the above example, scanner_name would be `github_repository_source_one` and scanner_type would be `github_repository`. This is 
important to filter resources in Grafana dashboard.