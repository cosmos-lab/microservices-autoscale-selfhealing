#!/bin/bash
set -e

echo ">>> Removing Helm releases..."
microk8s helm3 uninstall services || true
microk8s helm3 uninstall keda -n keda || true

echo ">>> Waiting for pods to terminate..."
microk8s kubectl delete pods --all --force --grace-period=0 2>/dev/null || true

echo ">>> Removing registry images..."
for SERVICE in order-service product-service inventory-service; do
  DIGEST=$(curl -s -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
    http://localhost:32000/v2/${SERVICE}/manifests/latest \
    -I | grep Docker-Content-Digest | awk '{print $2}' | tr -d '\r')

  if [ -n "$DIGEST" ]; then
    curl -s -X DELETE http://localhost:32000/v2/${SERVICE}/manifests/${DIGEST}
    echo "  Removed ${SERVICE}"
  fi
done

echo ""
echo ">>> Current state:"
microk8s kubectl get pods,svc,hpa 2>/dev/null || true

echo ""
echo ">>> Done. Run ./deploy.sh to bring everything back up."