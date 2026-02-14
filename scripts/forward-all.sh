#!/bin/sh

NAMESPACE="ecommerce"

# List of services and ports (multi-line)
SERVICES="
order-service:8080
payment-service:8081
inventory-service:8082
notification-service:8083
analytics-service:8084
read-service:8085
"

for svc_port in $SERVICES; do
  # Skip empty lines
  [ -z "$svc_port" ] && continue

  SERVICE=$(echo $svc_port | cut -d':' -f1)
  PORT=$(echo $svc_port | cut -d':' -f2)

  echo "Forwarding $SERVICE:$PORT ..."
  microk8s kubectl port-forward svc/$SERVICE $PORT:$PORT -n $NAMESPACE >/dev/null 2>&1 &
done

echo "All services are being port-forwarded in background."
echo "Use 'jobs' to see running port-forwards."
