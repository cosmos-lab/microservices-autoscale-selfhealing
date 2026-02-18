# Scalable Event-Driven Microservices on Kubernetes

## Tech Stack & Purpose

This project demonstrates how different cloud-native tools work together to build a scalable, production-style event-driven microservices system.


| Tool / Technology               | Purpose                                                               |
| ------------------------------- | --------------------------------------------------------------------- |
| Kubernetes (MicroK8s)           | Orchestrates containers, manages deployments, scaling and networking  |
| Helm                            | Templates and manages Kubernetes deployments consistently             |
| Apache Kafka                    | Enables async communication between microservices via event streaming |
| HPA (Horizontal Pod Autoscaler) | Scales HTTP-based services when CPU usage increases                   |
| KEDA                            | Scales event-driven services based on Kafka message lag               |
| NodePort Service                | Exposes service externally for local testing or API access            |
| Readiness Probe                 | Waits until pod is ready before allowing traffic routing              |
| Liveness Probe                  | Restarts pod automatically if application becomes unresponsive        |

---

## Autoscaling Strategy Used

Two different autoscaling techniques are used depending on workload type:

### CPU Based Scaling → HPA

Used for:

Order Service

Reason:

Order Service is user-facing and HTTP-driven.
Scaling based on CPU utilization ensures it handles increased incoming traffic.

---

### Event Driven Scaling → KEDA

Used for:

Inventory Service

Reason:

Inventory Service is a Kafka consumer.
Scaling based on:

Kafka Consumer Lag

ensures replicas increase when message backlog grows.

---

## Why Kafka?

Kafka enables asynchronous communication between services:

* Removes tight coupling between Order and Inventory
* Allows buffering of requests during traffic spikes
* Enables independent scaling of consumers
* Prevents cascading service failures

---

## Health Monitoring

Each service exposes:

```
/health
```

Used by Kubernetes for:

### Readiness Probe

* Prevents routing traffic to pods that are still starting

### Liveness Probe

* Automatically restarts pods stuck under heavy load

This prevents:

* Connection refused errors
* EOF failures
* Routing traffic to unhealthy containers

---

## Deployment Management

Helm is used for:

* Service deployments
* HPA creation
* KEDA ScaledObject creation
* Kafka configuration
* Kubernetes resource versioning

Ensures consistent and repeatable deployments across environments.

---

## Domain Design (DDD Inspired)

System is divided into independent business domains:

---

### Order Domain

Handled by:

Order Service

Responsibilities:

* Accept order request from client
* Validate product details
* Publish order event to Kafka topic `orders`

Scales using:

HPA (CPU-based autoscaling)

---

### Inventory Domain

Handled by:

Inventory Service

Responsibilities:

* Consume order events from Kafka
* Update product inventory asynchronously

Scales using:

KEDA (Kafka consumer lag)

---

## End-to-End Request Flow

1. Client sends POST request to Order Service
2. Order Service publishes event to Kafka topic `orders`
3. Inventory Service consumes event
4. Inventory updated asynchronously
5. Kafka lag increases under load
6. KEDA scales Inventory pods automatically
7. Lag decreases after processing

---

## Load Testing

Example:

```
hey -z 100s -c 50 -m POST \
-H "Content-Type: application/json" \
-d '{"productId":"1","productName":"Gaming Mouse","quantity":1}' \
http://localhost:30081/order
```

Expected Behavior:

* Order Service scales using HPA
* Inventory Service scales using KEDA
* No traffic routed to unhealthy pods
* Automatic recovery under load

---

## Scaling Summary

| Service           | Scaling Type | Trigger   |
| ----------------- | ------------ | --------- |
| Order Service     | HPA          | CPU Usage |
| Inventory Service | KEDA         | Kafka Lag |

---

This architecture demonstrates how synchronous APIs and asynchronous event consumers can be scaled independently in Kubernetes using HPA and KEDA.
