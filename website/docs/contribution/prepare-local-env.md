---
sidebar_position: 2
title: "Prepare Local Environment"
---

We have made it easy to get started working on TerediX codebase. You can follow the steps below to get started.

## Prerequisites
- Docker

## Clone the repository

```bash
git clone git@github.com:shaharia-lab/terediX.git
cd terediX
```

## Run TerediX inside Docker

We have a Docker image for development purpose. You can run the following command to start the development server.

```bash
docker run -i -d --name teredix-dev \
  -v $(pwd):/home/app/src \
  -p 3000:3000 \
  -p 2112:2112 \
  ghcr.io/shaharia-lab/teredix:dev
```

## Access the development environment

You can access development environment in Docker container by running the following command:

```bash
docker exec -it teredix-dev bash
```

Then, you can go to the project root

```bash
su app
cd ~/src/
```

## Test the development environment

```bash
make test-unit
```

## Run website

```bash
cd website
yarn install
yarn start --host=0.0.0.0 --port=3000
```

Now in your browser, you can access the website at http://localhost:3000

Voila! You are ready to contribute to TerediX.