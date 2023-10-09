---
sidebar_position: 3
title: "Getting Started"
---

# Getting Started

## Installation

There are several ways to install **terediX**. You can install **terediX** with standalone binary, `go install` command or with Docker.

### Standalone Binary

If you want to install **terediX** with binary, then you can download the binary from [release page](https://github.
com/shaharia-lab/teredix/releases). May be you need to make the binary executable by running the following command:

```bash
chmod +x terediX
```

### Go Install

If you want to install **terediX** with `go install` command, then you can run the following command:

```bash
go install github.com/shaharia-lab/teredix/cmd/terediX@latest
```

### Docker

If you want to install **terediX** with Docker, then you can run the following command to pull the docker image:

```bash
docker pull shaharialab/teredix:latest
```

Or, if you want to run the docker image, then you can run the following command:

```bash
docker run -it --rm shaharialab/teredix:latest --help
```

After installing the binary in your target machine, you can run the following command to see the help message:

```bash
teredix --help
```

## Create Configuration File

**terediX** uses a configuration file to run. You can create a configuration file with the following command:

```bash
teredix init
```

It will create a skeleton configuration file `config.yaml` in your current directory. You can edit the configuration 
file as per your need.

## Start Resource Scanner

After creating the configuration file, you can run the following command to run **terediX**:

```bash
teredix discover --config config.yaml
```