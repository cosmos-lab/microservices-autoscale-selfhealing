
# Infra

## Terraform
Infrastructure as Code (IaC) tool for provisioning cloud resources.

PLACE:
- Runs From:
    Local Machine / CI-CD Runner (GitHub Actions / Jenkins)
- Targets:
    Cloud Provider API (AWS / GCP / Azure)
- Creates:
    EKS (Elastic Kubernetes Service) Cluster
    EC2 (Elastic Compute Cloud) VM Nodes
    VPC (Virtual Private Cloud)
    Load Balancer (LB)

INPUT:
- Infrastructure Code (IaC)
- Desired Node Count
- Network Config

OUTPUT:
- Running Kubernetes Cluster Infra


## Helm
Kubernetes package manager for deploying applications and managing resources.

PLACE:
- Runs From:
    Local Machine / CI-CD Runner
    OR Inside Kubernetes via Helm Operator
- Deploys:
    Applications / Microservices
    Kubernetes Resources (Deployment, Service, ConfigMap, Secret, Ingress)

INPUT:
- Helm Chart (Package of Kubernetes Manifests + Templates)
- Values.yaml (Configuration Overrides)

OUTPUT:
- Installed / Upgraded Application in Kubernetes
- Sets initial desired state of pods, services, and configuration



## Kubernetes
Container orchestration platform that manages deployment, scaling, and self-healing.

PLACE:
- Lives On:
    EKS Control Plane (Managed by AWS)
    Worker Nodes = EC2 VM Instances

INPUT:
- Docker Container Images (from Helm Deploy)
- Deployment YAML (from Helm Chart)

Detects:
- Pod Crash (Container Process Exit)
- OOMKilled (Out Of Memory Killed)
- Readiness Probe Failure (Pod Not Ready To Serve Traffic)
- Liveness Probe Failure (Pod Not Alive)
- Node NotReady (Worker Node Down)

OUTPUT:
- Restart Container
- Recreate Pod
- Reschedule Pod

### HPA (Horizontal Pod Autoscaler)
Automatically scales Kubernetes pods based on observed metrics like CPU and memory.

PLACE:
- Built Inside Kubernetes Control Plane

INPUT:
- CPU / Memory Metrics From Prometheus

OUTPUT:
- Scale Pods 3 → 10 → 20

### KEDA (Kubernetes Event Driven Autoscaler)
Scales Kubernetes workloads based on external event sources like Kafka or Prometheus metrics. KEDA can scale pods based on queue length.

PLACE:
- Runs As:
    Kubernetes Operator Pod

INPUT:
- Queue Length
- Kafka Lag
- Custom Prometheus Metrics

OUTPUT:
- Scale Pods Based On Events

### Cluster Autoscaler
Automatically adjusts the number of nodes in a Kubernetes cluster based on pod resource requirements.

PLACE:
- Runs As:
    Kubernetes Pod

INPUT:
- Pending Pods
- Node Resource Exhaustion

OUTPUT:
- Requests AWS To Add EC2 Node

OpenTelemetry
Observability framework for collecting metrics, logs, and traces from applications.

PLACE:
- Runs As:
    DaemonSet (On Every Kubernetes Node)
    OR Sidecar Container With App Pod

INPUT:
- Pod Metrics
- Container Logs
- Network Latency
- API Request Duration
- DB Query Time
- Kubernetes Node & Pod Health

OUTPUT:
- Metrics → Prometheus
- Logs → ELK Stack
- Traces → Jaeger

### Jaeger
Distributed tracing system for monitoring and troubleshooting microservices.

PLACE:
- Runs As:
    Kubernetes Pod OR EC2

INPUT:
- Traces From OpenTelemetry

OUTPUT:
- Slow Service Detection
- DB Query Delay
- API Dependency Timeout

### LitmusChaos
Chaos engineering framework for testing system resilience and validating self-healing mechanisms.

PLACE:
- Runs As:
    Kubernetes Pod

INPUT:
- Chaos Experiment YAML

OUTPUT:
- Simulated Pod Kill
- Node Crash
- Network Delay
- Validates:
    Pod Recreation
    Auto Scaling
    Traffic Rerouting
    Deployment Rollback

### OPA (Open Policy Agent)
Policy enforcement tool for Kubernetes to ensure compliance and security.

PLACE:
- Runs As:
    Kubernetes Admission Controller

INPUT:
- Deployment Request

OUTPUT:
- Block Unsafe Deployment

### Istio
Service mesh that manages traffic routing, resilience, and security between services.

PLACE:
- Runs As:
    Sidecar Proxy Inside Each Pod

INPUT:
- Service To Service Traffic

OUTPUT:
- Circuit Breaker
- Retry Failed Request
- Route Traffic To Healthy Pod
- Outlier Detection (Kicks Failing Pods Out Before Crash)

        ↓


## Prometheus
Monitoring system and time-series database for metrics collection and alerting.

PLACE:
- Runs As:
    Kubernetes Pod (Monitoring Namespace)

INPUT:
- Metrics From OpenTelemetry

Detects:
- CPU Usage
- Memory Usage
- HTTP 5xx Error Rate
- Disk IO
- Pod Restart Count
- Network Packet Loss

OUTPUT:
- Sends Metrics To:
    HPA / KEDA
    Grafana Alert Rules

## ELK Stack
Centralized logging platform (Elasticsearch, Logstash, Kibana) for storing and visualizing logs.

PLACE:
- Elasticsearch → EC2 OR Kubernetes Pod
- Logstash → Kubernetes Pod
- Kibana → Kubernetes Pod OR EC2

INPUT:
- Application Logs
- Container Logs
- Exception Stacktrace
- API Timeout Logs
- DB Failure Logs

OUTPUT:
- Failure Pattern Detection
- Error Log Alerts
- Visualization via Kibana

## Grafana
Visualization and alerting platform for metrics and logs.

PLACE:
- Runs As:
    Kubernetes Pod OR EC2 Instance

INPUT:
- Metrics From Prometheus
- Logs From Elasticsearch

Detects:
- CPU > 70%
- Latency Increase
- Error Rate Spike
- Disk Usage > 90%
- DB Connection Limit Reached
- Thread Pool Exhaustion

OUTPUT:
- Sends Alerts To:
    HPA / KEDA
    Remediation Platform

## ArgoCD
GitOps continuous delivery tool for Kubernetes, keeping deployments in sync with Git repositories.

PLACE:
- Runs As:
    Kubernetes Pod

INPUT:
- Git Repo Desired State

Detects:
- CrashLoopBackOff
- Replica Mismatch
- Live State Drift (e.g., Manual CLI Changes to Replica Count)

OUTPUT:
- Rollback Deployment
- Heal Live State Drift

## PagerDuty Process Automation / Shoreline
Self-healing automation platform that executes remediation scripts based on alerts.

PLACE:
- SaaS Platform OR Agent Installed On EC2 / Kubernetes

INPUT:
- Alert From Grafana
- Metrics From Prometheus
- Logs From ELK

OUTPUT:
- Restart Service
- Clear Disk Cache
- Restart DB Pool
- Kill Hung Process
- Drain Node

---

### Explanation of Flow:

* **Terraform → Helm → Kubernetes → ArgoCD**
  Base deployment: IaC provisions infra, Helm deploys apps, Kubernetes orchestrates, ArgoCD ensures GitOps sync.

* **Kubernetes → HPA / KEDA → Kubernetes**
  Autoscaling triggered by metrics/events.

* **Kubernetes → Cluster Autoscaler → AWS Nodes**
  Scale nodes automatically based on pending pods or resource pressure.

* **Kubernetes → OpenTelemetry → Prometheus / ELK Stack → Grafana → HPA / KEDA**
  Observability loop: metrics/logs feed auto-scaling and alerting.

* **Grafana → PagerDuty / Shoreline**
  Alerts trigger self-healing actions.

* **Kubernetes → Istio**
  Service-to-service routing, retries, circuit breaking, outlier detection.

* **Kubernetes → Jaeger**
  Tracing for slow services, DB delays, and dependency timeouts.

* **Kubernetes → OPA**
  Policy enforcement to block unsafe deployments.

* **LitmusChaos → Kubernetes**
  Chaos experiments simulate failures to test auto-healing.

* **PagerDuty / Shoreline → Kubernetes**
  Automatic remediation scripts act on detected failures.


## Project Structure

```
/root
│
├── terraform/               # INFRASTRUCTURE LAYER (Provider: AWS/Azure/GCP)
│   ├── main.tf              # VPC, EKS/GKE Cluster, IAM Roles
│   ├── variables.tf
│   ├── outputs.tf
│   └── modules/             # Reusable Terraform modules (VPC, EKS, RDS)
│
├── helm-charts/             # PACKAGING LAYER (Application Blueprints)
│   ├── api-app/             # Custom Application Chart
│   │   ├── Chart.yaml
│   │   ├── values.yaml      # Default Dev values
│   │   └── templates/       # Deployment, Service, HPA, PodMonitor
│   │
│   ├── database/            # Custom Postgres/Redis Chart
│   │   ├── Chart.yaml
│   │   └── templates/       # StatefulSet, PVC, Secrets
│   │
│   ├── observability/       # Helm charts for Prometheus, Grafana, OpenTelemetry
│   │   ├── prometheus/      # Prometheus Helm chart + values.yaml
│   │   ├── grafana/         # Grafana Helm chart + dashboards
│   │   └── otel-collector/  # OpenTelemetry Collector Helm chart + values
│   │
│   ├── istio/               # Istio Helm or manifests
│   │   └── values.yaml
│   │
│   ├── opa/                 # OPA Gatekeeper / Helm chart values
│   │   └── values.yaml
│   │
│   └── keda/                # KEDA Helm chart values
│       └── values.yaml
│
├── observability/           # CONFIG + DASHBOARDS
│   ├── grafana-dashboards/       # JSON files for dashboards
│   ├── prometheus-values.yaml    # Custom overrides for Prometheus
│   └── otel-collector-values.yaml# Custom OpenTelemetry Collector config
│
├── gitops/                  # DEPLOYMENT STATE LAYER (ArgoCD)
│   ├── base/                # Common config for all environments
│   │   ├── apps.yaml        # Base applications manifests / HelmRelease references
│   │   └── namespaces.yaml
│   │
│   └── overlays/            # Environment-specific overrides
│       ├── production/
│       │   ├── kustomization.yaml
│       │   └── values-prod.yaml  # Prod-specific Helm values
│       │
│       └── staging/
│           ├── kustomization.yaml
│           └── values-staging.yaml  # Staging-specific Helm values
│
├── microservices/           # NODE.JS MICROSERVICES
│   ├── product/
│   │   ├── cmd/server.js
│   │   ├── proto/product.proto
│   │   ├── internal/service.js
│   │   ├── Dockerfile
│   │   └── helm/
│   │
│   ├── order/
│   │   ├── cmd/server.js
│   │   ├── proto/order.proto
│   │   ├── internal/service.js
│   │   ├── kafka/producer.js
│   │   ├── kafka/consumer.js
│   │   ├── Dockerfile
│   │   └── helm/
│   │
│   ├── inventory/
│   ├── payment/
│   ├── shipping/
│   ├── notification/
│   ├── user/
│   └── shared/
│       ├── kafka-client.js
│       ├── redis-client.js
│       └── grpc-client.js
│
├── chaos/                   # Chaos Engineering Experiments (LitmusChaos)
│   ├── pod-delete.yaml
│   ├── network-delay.yaml
│   └── cpu-stress.yaml
│
└── remediation/             # Self-healing scripts / runbooks (PagerDuty/Shoreline)
    ├── clear-disk-cache.sh
    ├── restart-db.sh
    └── flush-tcp-buffers.sh

```
