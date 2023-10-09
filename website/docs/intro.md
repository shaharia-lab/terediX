---
sidebar_position: 1
title: "terediX Intro"
---

# Introduction

`TeReDiX` (Tech Resource Discover &amp; Exploration) is a tool to discover tech resources for an organization from 
various resource provider (e.g. AWS, GitHub, GCP, Atlassian, etc.).

**terediX** can be a useful tool for organizations who have so many resources from different providers, and they want 
to track, monitor those resources for better visibility across the organization. This tool can efficiently fetch all 
the resources and it's associated metadata from different providers and store them in a central database. Later, the 
data can be visualize in a dashboard for better visibility. 

TLDR; **terediX** is a tool to manage your tech resource inventory centrally.

<p align="center">
  <a href="https://github.com/shaharia-lab/teredix"><img src="https://user-images.githubusercontent.com/1095008/229536376-51ddaa75-85ee-4e3e-95df-7cf6093392bf.png" width="100%"/></a>
</p><br/>

## Features
Here is the list of top features currently **terediX** can offer:

- [x] Fetch resources from supported providers
- [x] Store resources in a central database
- [x] Expose prometheus metrics
- [x] Visualize resources in Grafana dashboard
- [x] Support multiple providers
- [x] Fully configurable with a single configuration YAML file
- [x] Support multiple database backends
- [x] Support customized metadata for each resource type
- [x] Support multiple resource types (e.g: EC2, S3, RDS, File system, GitHub repository, etc.)
- [x] Support multiple deployment options (e.g: Single binary, Docker, Kubernetes with Helm, etc.)
- [x] Optimized for large scale resource discovery
- [x] Customized scheduler for each resource scanner
- [x] Support re-usable configuration among multiple resource scanner
- [x] Support multiple resource scanner for same resource type
- [x] Fully open source, free to use & community driven initiative

### Technical Architecture

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

The scheduler is responsible for scheduling the execution of the resource scanners. The resource scanners are responsible for fetching resources from their respective sources and building metadata for each resource. The processor is responsible for processing the resources and storing them in the storage.

The architecture is designed to be scalable and extensible. The scheduler can be configured to schedule the resource scanners to run in parallel, and the processor can be configured to process resources in batches. The storage can be implemented using a variety of technologies, such as a database or a distributed file system.

Here is a more detailed explanation of each component:

**Scheduler:** The scheduler is responsible for scheduling the execution of the resource scanners. It can be configured 
to schedule the scanners to run in parallel, and it can also be configured to schedule them to run on different machines. 
The scheduler can also be configured to run the scanners at specific times, or at regular intervals.

**Resource Scanner:** The resource scanners are responsible for fetching resources from their respective sources and building metadata for each resource. 
The resource scanners can be implemented using a variety of technologies, such as APIs or web scraping.

**Processor:** The processor is responsible for processing the resources and storing them in the storage. 
The processor can be configured to perform a variety of tasks on the resources, such as filtering, transforming, and enriching the data. The processor can also be configured to store the resources in different formats.

**Storage:** The storage is responsible for storing the processed resources. The storage can be implemented using a variety of technologies, 
such as a database or a distributed file system.

The **terediX** architecture is designed to be flexible and adaptable. It can be used to process a wide variety of resources from a wide variety of sources. 
The architecture can also be scaled to meet the needs of small and large organizations.

## Generate a new site

Generate a new Docusaurus site using the **classic template**.

The classic template will automatically be added to your project after you run the command:

```bash
npm init docusaurus@latest my-website classic
```

You can type this command into Command Prompt, Powershell, Terminal, or any other integrated terminal of your code editor.

The command also installs all necessary dependencies you need to run Docusaurus.

## Start your site

Run the development server:

```bash
cd my-website
npm run start
```

The `cd` command changes the directory you're working with. In order to work with your newly created Docusaurus site, you'll need to navigate the terminal there.

The `npm run start` command builds your website locally and serves it through a development server, ready for you to view at http://localhost:3000/.

Open `docs/intro.md` (this page) and edit some lines: the site **reloads automatically** and displays your changes.
