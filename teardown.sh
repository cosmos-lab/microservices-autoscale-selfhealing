#!/bin/bash
set -e

REGISTRY="localhost:32000"

echo ">>> Uninstall Helm Releases..."
microk8s helm3 uninstall services || true
microk8s helm3 uninstall keda -n keda || true

echo ">>> Delete leftover workloads..."
microk8s kubectl delete deploy,rs,po,svc,hpa,scaledobject --all --grace-period=0 --force || true

echo ">>> Delete Kafka PVC junk (offsets)..."
microk8s kubectl delete pvc --all || true

echo ">>> Remove KEDA CRDs..."
microk8s kubectl delete crd scaledobjects.keda.sh 2>/dev/null || true
microk8s kubectl delete crd scaledjobs.keda.sh 2>/dev/null || true
microk8s kubectl delete crd triggerauthentications.keda.sh 2>/dev/null || true
microk8s kubectl delete crd clustertriggerauthentications.keda.sh 2>/dev/null || true

echo ">>> Removing MicroK8s containerd cached images..."

for SERVICE in order-service product-service inventory-service; do
microk8s ctr images rm ${REGISTRY}/${SERVICE}:latest 2>/dev/null || true
done

echo ">>> Remove Kafka/Zookeeper images too..."
microk8s ctr images rm confluentinc/cp-kafka:7.5.0 2>/dev/null || true
microk8s ctr images rm confluentinc/cp-zookeeper:7.5.0 2>/dev/null || true

echo ">>> Pruning ALL unused containerd layers..."
microk8s ctr content rm $(microk8s ctr content ls -q) 2>/dev/null || true

echo ">>> Restart container runtime..."
microk8s stop
sleep 5
microk8s start

echo ""
echo ">>> CLUSTER IS NOW CLEAN"
echo "Run: ./deploy.sh"
