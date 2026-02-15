

# Test APIs

```
curl -X POST http://localhost:8080/order/create \
     -H "Content-Type: application/json" \
     -d '{"productId": "p-99", "productName": "Gaming Mouse"}'

//Load test order
hey -z 10s -c 50 -m POST \
    -H "Content-Type: application/json" \
    -d '{"productId": "p-99", "productName": "Gaming Mouse"}' \
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
"./tmp" \
"./main.go" \
"./tmp" \
"./services/product/go.mod,"\
"./services/product/go.sum,"\
"./services/product/go.mod,"\
"./services/product/go.sum" \
"output.txt"


../tools/llmcode/export_code.sh \
"./helm,"\
"./tmp" \
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

---

**Pods — how many are running**
```bash
microk8s kubectl get pods
microk8s kubectl get pods -w
```

---

**HPA — current replicas + CPU vs target**
```bash
microk8s kubectl get hpa
# detailed view with events (scale up/down history)
microk8s kubectl describe hpa order-service
microk8s kubectl describe hpa product-service
```
The `get hpa` output is the most useful — it shows:
```
NAME              MINPODS  MAXPODS  REPLICAS  TARGETS    AGE
order-service     1        5        1         12%/60%    2m
product-service   1        5        1         8%/60%     2m
```
`TARGETS` = `current%/threshold%` — when current hits 60% it will scale up.

---

**Metrics — actual CPU/memory per pod**
```bash
microk8s kubectl top pods
microk8s kubectl top nodes
```

---

**Watch everything at once (best for load testing)**
```bash
watch -n 0.5 'microk8s kubectl get pods,hpa && echo "" && microk8s kubectl top pods'
```
This refreshes every 2 seconds showing pods, HPA state, and live CPU together in one terminal.

---

**Trigger load to actually see scaling** — simplest way with just `curl` in a loop:
```bash
for i in $(seq 1 1000); do
  curl -s -X POST http://localhost:30081/order/create \
    -H "Content-Type: application/json" \
    -d '{"productId":"1","productName":"Laptop"}' &
done
```
Then watch the `watch` command above — you should see `REPLICAS` climb from 1 toward 5 as CPU crosses 60%.