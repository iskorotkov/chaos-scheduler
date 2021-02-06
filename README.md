# Chaos Scheduler

Service for automatic generation and scheduling of the chaos test workflows.

- [Chaos Scheduler](#chaos-scheduler)
  - [Overview](#overview)
    - [Failures](#failures)
      - [Container failures](#container-failures)
      - [Pod, part of deployment, deployment failures](#pod-part-of-deployment-deployment-failures)
      - [Node failures](#node-failures)
      - [Cluster failures](#cluster-failures)
    - [Targets](#targets)
    - [Workflows](#workflows)
  - [Setup](#setup)
    - [Dependencies](#dependencies)
    - [Installation](#installation)
    - [Env var](#env-var)
  - [Extension](#extension)
    - [REST API](#rest-api)
    - [Annotations](#annotations)
    - [Development](#development)
    - [Project structure](#project-structure)

## Overview

Service generates workflows consisting of failures. Each failure has associated target. Each failure is something bad happening to the target (network loss, pod deletion, etc). Target is a part of the system under test (container, pod, deployment, etc). Service fetches list of potential targets from Kubernetes.

### Failures

**Failures selection**. Service picks failures in random order according to several constraints: max number of failures per stage, max chaos score per stage. It's possible to change failures selection by providing another seed.

**Failure score**. Each failure has chaos score - amount of chaos that it will create once executed. Failure score depends on its scale and severity.

#### Container failures

Container failures target a single container in a pod:

- container CPU hog
- container memory hog
- container network corruption
- container network duplication
- container network latency
- container network loss

#### Pod, part of deployment, deployment failures

Pod failures target pods, parts of deployments (specified percent of pods in deployment) or entire deployments:

- pod delete
- pod I/O stress

#### Node failures

Node failures target entire nodes:

- node CPU hog
- node memory hog
- node I/O stress

#### Cluster failures

Cluster failures target entire cluster. No cluster failures currently supported.

### Targets

**Namespace**. Targets are fetched from Kubernetes. Service fetches only targets from specified namespace to avoid targeting production instances.

**Annotation**. Moreover, each target must have annotation `litmuschaos.io/chaos: "true"` in order for chaos experiments to work. It is a precaution measure to avoid targeting production instances.

**Target selection**. Service randomly picks appropriate target for each failure. It's possible to change target selection by providing another seed.

### Workflows

**Structure**. Each workflow consists of specified number of test stages. Each stage consists of several steps (or actions). Each step is a failure with associated target.

**Execution order**. Stages are executed in order one at a time. Steps in each stage are executed at the same time.

**Workflow preview**. Service allows previewing generated workflow without launching it.

## Setup

### Dependencies

Install dependencies before continuing:

- Argo
- Litmus Chaos

You also have to create `ServiceAccount` for Litmus Chaos and Argo.

### Installation

Make sure you have a Kubernetes cluster ready.

Install all dependencies and make sure they work correctly.

Tweak env var values in `deploy/scheduler.yaml` file for your environment (optional). Once finished, execute the command in the root folder:

```shell
kubectl -f deploy/scheduler.yaml
```

### Env var

Service requires several env vars set (example values are provided in parentheses):

- ARGO_SERVER - Argo server to use (`argo-server.argo.svc:2746`)
- STAGE_MONITOR_IMAGE - Docker image to use for monitoring crashes of target containers/pods (`iskorotkov/chaos-pods-monitor:v0.4.0`)
- APP_NS - namespace where system under test is located (`chaos-app`)
- CHAOS_NS - namespace where to create workflows (`litmus`)
- APP_LABEL - label to use for target selection (`app`)

    Service looks for label `{APP_LABEL}: {VALUE}`, where `{VALUE}` will be the name of the target.

    For example, when `APP_LABEL`=`app` the service will look for label `app: {VALUE}`. The target with label `app: nginx` will be named `nginx`.

- DEVELOPMENT - whether in development or not (false)
- STAGE_DURATION - duration of each stage (30s)
- STAGE_INTERVAL - duration between stages (30s)

    Some failures take seconds to start and can't finish instantly. It's recommended to set interval to 30s or higher to avoid false positives in latter stages.

## Extension

### REST API

- api/v1/workflows
  - preview - generate and preview test workflow (without launching it)
  - create - generate and launch test workflow

### Annotations

Service adds annotations to generated workflow steps according:

| Key                         | Category       | Type   | Description                                                                |
| --------------------------- | -------------- | ------ | -------------------------------------------------------------------------- |
| chaosframework.com/version  | version        | semver | Version of annotations format                                              |
| chaosframework.com/type     | classification | string | Type of template (failure, utility)                                        |
| chaosframework.com/severity | classification | string | Failure severity (harmless, light, severe, critical)                       |
| chaosframework.com/scale    | classification | string | Failure scale (container, pod, deployment part, deployment, node, cluster) |

### Development

To build project:

```shell
go build ./...
```

To run tests:

```shell
go test ./...
```

### Project structure

- cmd
  - scheduler - entry point
- internal
  - handlers - request handlers
  - config - getting config from environment
- pkg
  - argo - argo client for executing workflows
  - k8s - kubernetes client for fetching list of targets
  - rx - random string, map and slice generation
  - server - advanced request handling
  - workflows - workflow creation and execution
    - generate - test scenario creation

      Scenario is a logical representation of a chaos test, while workflow is a 1) practical representation of a chaos test; 2) scenario prepared to be executed.

    - assemble - test workflow creation
    - execution - test workflow execution
