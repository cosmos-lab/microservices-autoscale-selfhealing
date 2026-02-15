#!/bin/bash
set -e

REGISTRY="localhost:32000"

echo ">>> Enabling microk8s addons..."
microk8s enable registry helm3 metrics-server

echo ">>> Installing KEDA..."
microk8s helm3 repo add kedacore https://kedacore.github.io/charts 2>/dev/null || true
microk8s helm3 repo update
microk8s helm3 upgrade --install keda kedacore/keda \
  --namespace keda --create-namespace \
  --version 2.13.0

echo ">>> Waiting for KEDA CRDs..."
until microk8s kubectl explain scaledobject 2>/dev/null | grep -q "KIND"; do
  echo "  not ready, retrying in 3s..."
  sleep 3
done
echo "  KEDA CRDs ready!"

echo ">>> Building images with Podman..."
podman build -t ${REGISTRY}/order-service:latest     ./services/order
podman build -t ${REGISTRY}/product-service:latest   ./services/product
podman build -t ${REGISTRY}/inventory-service:latest ./services/inventory

echo ">>> Pushing to microk8s registry..."
podman push ${REGISTRY}/order-service:latest     --tls-verify=false
podman push ${REGISTRY}/product-service:latest   --tls-verify=false
podman push ${REGISTRY}/inventory-service:latest --tls-verify=false

echo ">>> Deploying services (without KEDA)..."
microk8s helm3 upgrade --install services ./helm/services \
  --set inventory-service.keda=null

echo ">>> Waiting for services to be ready..."
microk8s kubectl rollout status deployment/order-service --timeout=120s
microk8s kubectl rollout status deployment/product-service --timeout=120s
microk8s kubectl rollout status deployment/inventory-service --timeout=120s

echo ">>> Applying KEDA ScaledObjects..."
microk8s helm3 upgrade --install services ./helm/services

echo ""
echo ">>> Status:"
microk8s kubectl get pods,svc,hpa

echo ""
echo ">>> KEDA ScaledObjects:"
microk8s kubectl get scaledobjects

echo ""
echo ">>> Endpoints:"
echo "  Order:     http://localhost:30081/order/create"
echo "  Product:   http://localhost:30082/product/:id"
echo "  Inventory: http://localhost:30083/inventory/:productId"