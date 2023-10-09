---
title: "File System"
---

# File System

<img src="/img/file_system_icon.png" alt="File System" width="250"/>

## Configuration

### Type

Resource type. In this case it would be `file_system`.

### Configuration

- **root_directory**: Provide the absolute path of the directory to scan.

### Fields

List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data 
from the following fields.

- machineHost
- rootDirectory

### Schedule

**🔗 [Check schedule format](/docs/configuration/scanner/overview#schedule-format)**.

## Example

```yaml
source:
  file_system_source_one:
      type: file_system
      configuration:
        root_directory: "/file/path"
      fields:
        - machineHost
        - rootDirectory
      schedule: "@every 24h"
```

In the above example, we have added a source named `file_system_source_one` with type `file_system`. We have added some fields to add with each resource. 
We have also set the schedule to run this source every 24 hours.

Based on the above example, scanner_name would be `file_system_source_one` and scanner_type would be `file_system`. This is 
important to filter resources in Grafana dashboard.