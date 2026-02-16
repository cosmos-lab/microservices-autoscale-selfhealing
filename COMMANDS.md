

# Test APIs

```
curl -X POST http://localhost:8080/order/create \
     -H "Content-Type: application/json" \
     -d '{"productId": "p-99", "productName": "Gaming Mouse", "quantity": 1}'

//Load test order
hey -z 100s -c 50 -m POST \
    -H "Content-Type: application/json" \
    -d '{"productId": "1", "productName": "Gaming Mouse", "quantity": 1}' \
    http://localhost:30081/order/create

hey -z 120s -c 200 -m POST \
    -H "Content-Type: application/json" \
    -d '{"productId": "1", "productName": "Gaming Mouse", "quantity": 1}' \
    http://localhost:30081/order/create


curl -X POST http://localhost:8080/product \
     -H "Content-Type: application/json" \
     -d '{"id": "101", "name": "Mechanical Keyboard"}'


//Load test product
hey -z 10s -c 50 -m POST \
    -H "Content-Type: application/json" \
    -d '{"id": "test", "name": "stress-test"}' \
    http://localhost:30082/product




```

# **LLM Code** #

```
../tools/llmcode/export_code.sh \
"./services,"\
"./helm,"\
"./tmp" \
"./main.go" \
"./tmp" \
"./services/product/go.mod,"\
"./services/product/go.sum,"\
"./services/order/go.mod,"\
"./services/order/go.sum" \
"output.txt"


../tools/llmcode/export_code.sh \
"./helm" \
"./tmp.go" \
"./tmp" \
"./tmp.go" \
"output.txt"


../tools/llmcode/export_code.sh \
"./services"\
"./tmp.go" \
"./tmp" \
"./tmp.go" \
"output.txt"

../../../tools/llmcode/export_structure.sh \
  ./ \
  tmp \
  output.txt


../tools/llmcode/import_code.sh input.txt ./

```

# Monitor

```bash
watch -n 1 '
echo "=== PODS ==="
microk8s kubectl get pods

echo ""
echo "=== HPA ==="
microk8s kubectl get hpa

echo ""
echo "=== KEDA ==="
microk8s kubectl get scaledobjects

echo ""
echo "=== CPU/MEM ==="
microk8s kubectl top pods

echo ""
echo "=== KAFKA LAG ==="
microk8s kubectl exec -it \
  $(microk8s kubectl get pod -l app=kafka -o jsonpath="{.items[0].metadata.name}") \
  -- kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe \
  --group inventory-service 2>/dev/null
'
```
---


# Logs
```
# logs
microk8s kubectl logs -l app=kafka

# logs for init container
microk8s kubectl logs -l app=kafka -c wait-for-zookeeper

# describe
microk8s kubectl describe pod -l app=kafka

# exec into it
microk8s kubectl exec -it -l app=kafka -- sh

microk8s kubectl logs -l app=kafka --previous


```


#

```
microk8s kubectl get all -n kube-system

``