#!/bin/bash
set -e

REGISTRY="localhost:32000"

echo ">>> Installing KEDA..."
microk8s helm3 repo add kedacore https://kedacore.github.io/charts 2>/dev/null || true
microk8s helm3 repo update
microk8s helm3 upgrade --install keda kedacore/keda \
  --namespace keda --create-namespace \
  --version 2.13.0

echo ">>> Waiting for KEDA CRDs..."
until microk8s kubectl explain scaledobject 2>/dev/null | grep -q "KIND"; do
  sleep 3
done

echo ">>> Waiting for Metrics Server..."
until microk8s kubectl get apiservice v1beta1.metrics.k8s.io 2>/dev/null | grep -q True; do
  sleep 3
done

echo ">>> Building images..."
podman build -t ${REGISTRY}/order-service:latest     ./services/order
podman build -t ${REGISTRY}/product-service:latest   ./services/product
podman build -t ${REGISTRY}/inventory-service:latest ./services/inventory

echo ">>> Pushing images..."
podman push ${REGISTRY}/order-service:latest     --tls-verify=false
podman push ${REGISTRY}/product-service:latest   --tls-verify=false
podman push ${REGISTRY}/inventory-service:latest --tls-verify=false

echo ">>> Deploying services..."
microk8s helm3 upgrade --install services ./helm/services \
  --reset-values \
  --atomic \
  --timeout 5m

echo ">>> Waiting for Kafka..."
microk8s kubectl rollout status deployment/kafka --timeout=180s

echo ">>> Creating Kafka Topics..."

create_topic () {
microk8s kubectl exec deployment/kafka -- \
kafka-topics \
--create \
--if-not-exists \
--topic $1 \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1
}

create_topic orders
create_topic inventory-reserved
create_topic inventory-failed

echo ">>> Restarting consumers after topic creation..."
microk8s kubectl rollout restart deployment inventory-service
microk8s kubectl rollout restart deployment order-service

echo ">>> Waiting for services..."
microk8s kubectl rollout status deployment/order-service --timeout=120s
microk8s kubectl rollout status deployment/product-service --timeout=120s
microk8s kubectl rollout status deployment/inventory-service --timeout=120s

echo ""
echo ">>> Status:"
microk8s kubectl get pods,svc,hpa
