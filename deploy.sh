#!/bin/bash
set -e

REGISTRY="localhost:32000"

echo ">>> Enabling microk8s addons..."
microk8s enable registry helm3 metrics-server

echo ">>> Building images with Podman..."
podman build -t ${REGISTRY}/order-service:latest   ./services/order
podman build -t ${REGISTRY}/product-service:latest ./services/product

echo ">>> Pushing to microk8s registry..."
podman push ${REGISTRY}/order-service:latest   --tls-verify=false
podman push ${REGISTRY}/product-service:latest --tls-verify=false

echo ">>> Deploying with Helm..."
microk8s helm3 upgrade --install services ./helm/services

echo ""
echo ">>> Status:"
microk8s kubectl get pods,svc,hpa

echo ""
echo ">>> Endpoints:"
echo "  Order:   http://localhost:30081/order/create"
echo "  Product: http://localhost:30082/product/:id"
