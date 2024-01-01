# terediX Helm Chart

[terediX](https://github.com/shaharia-lab/terediX) (Tech Resource Discovery & Exploration) is a tool to discover tech resources for an organization and explore them

## Introduction

This helm chart can be used to deploy [terediX](https://github.com/shaharia-lab/terediX) in Kubernetes using [Helm](https://helm.sh) package package manager

## Prerequisites

- Helm 3+
- Kubernetes 1.12+

## Installing the Chart

To add the repository to Helm:

```shell
helm repo add teredix https://shaharia-lab.github.io/terediX
```

## Install locally

First, you need to create a namespace for `terediX`, then need to install PostgreSQL and then you can install the application

```bash
kubectl create ns tererdix

helm repo add https://charts.bitnami.com/bitnami
helm repo update

helm upgrade --install postgresql bitnami/postgresql --namespace "teredix" \
        --set auth.username="app" \
        --set auth.password="pass" \
        --set auth.database="app"
```

Build Docker image

```bash
docker build -t teredix:prod -f Dockerfile .

helm upgrade --install teredix ./helm-chart/teredix --namespace teredix \
        -f ./helm-chart/teredix/values.yaml \
        -f ./helm-chart/teredix/values-local.yaml
```
