---
sidebar_position: 2
title: "Technical Architecture"
---

# Technical Architecture

```
                                    +----------+
                                    | Scheduler |
                                    +----------+
                                          |
                                          |
                                          v
                +--------------+     +--------------+
                | Resource     |     | Resource     |
                | Scanner 1    |     | Scanner 2    |
                +--------------+     +--------------+
                                          |
                                          |
                                          v
                                  +----------+
                                  | Processor |
                                  +----------+
                                          |
                                          |
                                          v
                                  +----------+
                                  | Storage   |
                                  +----------+
```
## Components

### Configuration file

The configuration file is a JSON file that contains all the necessary settings for TerediX to run. The configuration file specifies the following:

- The list of scanners to use.
- The schedule for each scanner.
- The storage configuration.
- The metrics configuration.

### Scanner

Scanners are responsible for fetching resources from their respective sources. Scanners can be implemented in any programming language. 
TerediX includes a number of built-in scanners, such as a scanner for AWS S3 and a scanner for Google Cloud Storage.

### Scheduler

The scheduler is responsible for scheduling scanners to run at regular intervals. The scheduler can be implemented 
in any programming language. TerediX includes a built-in scheduler that uses cron expressions.

### Processor

The processor is responsible for processing resources as they are fetched by scanners. The processor can be implemented in any programming language.
TerediX includes a built-in processor that stores resources and metadata in the storage.

### Storage

The storage is responsible for storing resources and metadata. The storage can be implemented in any database or file system.
TerediX includes a built-in storage that uses a PostgreSQL database.

### Metrics

The metrics component collects and exposes metrics about TerediX's operation. The metrics component can be implemented in any programming language.
TerediX includes a built-in metrics component that exposes metrics to Prometheus.

**The components are connected as follows:**

- The processor reads the configuration file to determine which scanners to use and their schedules.
- The processor starts the scanners and processes the resources that they fetch.
- The processor stores the resources and metadata in the storage.
- The metrics component collects and exposes metrics about the processor's operation.

The system is designed to be scalable and reliable. The processor can be scaled horizontally to handle more resources. The storage is designed to be highly available and durable.
