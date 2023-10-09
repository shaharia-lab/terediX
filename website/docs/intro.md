---
sidebar_position: 1
title: "Introduction"
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
