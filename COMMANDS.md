

# Test APIs

```
curl -X POST http://localhost:8080/order/create \
     -H "Content-Type: application/json" \
     -d '{"productId": "p-99", "productName": "Gaming Mouse"}'


curl -X POST http://localhost:8080/product \
     -H "Content-Type: application/json" \
     -d '{"id": "101", "name": "Mechanical Keyboard"}'


//Load test product
hey -z 10s -c 50 -m POST \
    -H "Content-Type: application/json" \
    -d '{"id": "test", "name": "stress-test"}' \
    http://localhost:8080/product


//Load test order
hey -z 10s -c 50 -m POST \
    -H "Content-Type: application/json" \
    -d '{"productId": "p-99", "productName": "Gaming Mouse"}' \
    http://localhost:8080/order/create


```

# Autoscale

```
microk8s enable metrics-server
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

../../../tools/llmcode/export_structure.sh \
  ./ \
  tmp \
  output.txt


../tools/llmcode/import_code.sh input.txt ./

```
