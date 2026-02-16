set -e

echo ">>> Enabling microk8s addons..."
microk8s enable registry helm3 metrics-server dns

echo ">>> Waiting for CoreDNS pod to be ready..."
microk8s kubectl rollout status deployment/coredns -n kube-system --timeout=60s

echo ">>> Ensuring CoreDNS service exists..."
microk8s kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    k8s-app: kube-dns
spec:
  selector:
    k8s-app: kube-dns
  clusterIP: 10.152.183.10
  ports:
    - name: dns
      port: 53
      protocol: UDP
    - name: dns-tcp
      port: 53
      protocol: TCP
EOF

echo ">>> Verifying DNS..."
microk8s kubectl get svc -n kube-system

echo ""
echo ">>> Done. Run ./deploy.sh to deploy services."