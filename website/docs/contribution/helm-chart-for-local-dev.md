---
sidebar_position: 3
title: "Helm Chart for Local Development"
---

We have made it easy to test deploying terediX in your local Kubernetes cluster. You can follow the steps below to get started.

## Prerequisites

- Docker
- Kubernetes Cluster (Minikube, Kind, K3s, etc). You can easily start a local Kubernetes cluster using [KiND](https://github.com/shaharia-lab/k8s-dev-cluster)

Go to the next step when your Kubernetes cluster is ready.

## Create a namespace

Create a namespace for terediX

```bash
kubectl create namespace teredix
```

## Install PostgreSQL

Because terediX need a storage solution and currently [we support](/docs/configuration/storage#supported-storage-engines) only PostgreSQL, 
so you need to install PostgreSQL in your Kubernetes cluster. You can install PostgreSQL by helm chart.

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami 
helm repo update
helm upgrade --install postgresql bitnami/postgresql --namespace "teredix" \
        --set auth.username="app" \
        --set auth.password="pass" \
        --set auth.database="app"
```

## Install terediX helm chart

Create a local values file in `helm-chart/teredix/values-local.yaml` for terediX helm chart to override few values for local development.

```yaml
# helm-chart/teredix/values-local.yaml
image:
  repository: teredix
  tag: "prod"

appConfig:
  organization:
    name: Your Organization
    logo: https://your-org-url.com/logo.png
  discovery:
    name: Name of the discovery
    description: Some description about the discovery
    worker_pool_size: 1
  storage:
    batch_size: 2
    engines:
      postgresql:
        host: "postgresql"
        port: 5432
        user: "app"
        password: "pass"
        db: "app"
    default_engine: postgresql
  source:
    fs_one:
      type: file_system
      configuration:
        root_directory: "/config"
      fields:
        - machineHost
        - rootDirectory
      schedule: "@every 300s"
  relations:
    criteria:
      - name: "file-system-rule1"
        source:
          kind: "FilePath"
          meta_key: "rootDirectory"
          meta_value: "/some/path"
        target:
          kind: "FilePath"
          meta_key: "rootDirectory"
          meta_value: "/some/path"

service:
  type: ClusterIP
  port: 2112

ingress:
  enabled: true
  hosts:
    - host: teredix.dev.local
      paths:
        - path: /
          pathType: ImplementationSpecific
```

Now install terediX helm chart using the following command:

```bash
helm upgrade --install teredix ./helm-chart/teredix --namespace teredix \
        -f ./helm-chart/teredix/values.yaml \
        -f ./helm-chart/teredix/values-local.yaml
```