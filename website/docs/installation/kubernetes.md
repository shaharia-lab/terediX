---
title: "Deploy terediX in Kubernetes"
sidebar_label: Kubernetes
---

## Deploy in Kubernetes using Helm Chart

It's very simple to install `terediX` in Kubernetes. You can install using Helm chart. Here are the simplest steps to follow.

### Add Helm Repository

```bash
helm repo add teredix https://teredix.shaharialab.com
helm repo update
```

### Add Configuration

Create a separate helm values file and override necessary configuration. Specially you need to provide the configuration for `terediX`.

Create a `values-prod.yaml` file and put the following content. You can override any configuration as you need. Read more about [terediX configuration](/docs/configuration/general).

```bash
appConfig:
#  organization:
#     name: Your Organization
#     logo: https://your-org-url.com/logo.png
#   discovery:
#     name: Name of the discovery
#     description: Some description about the discovery
#     worker_pool_size: 1
#   storage:
.........
```

### Install
After that, just install terediX using the following command.

```bash
helm install teredix teredix/teredix --namespace teredix -f values-prod.yaml
```

For more useful Helm commands, please follow the [official documentation of Helm](https://helm.sh/docs/helm/).

For details about `terediX` helm chart, go to [terediX Helm Chart in ArtifactHub](https://artifacthub.io/packages/helm/teredix/teredix)