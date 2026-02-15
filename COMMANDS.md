

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


## Deploy all microservices via Helm 


### dev mode
```


microk8s helm3 install microservices ./helm/microservices \
  --namespace ecommerce \
  --create-namespace


#Check pods
microk8s kubectl get pods -n ecommerce

#Check Internal Services (Cluster Networking)
microk8s kubectl get svc -n ecommerce


#Check logs
microk8s kubectl describe order-service -n ecommerce


#Get IP:
microk8s kubectl get nodes -o wide

http://<NODE-IP>:30080/order/create

```

### dev prod
```
microk8s helm3 upgrade --install microservices ./helm-charts/microservices \
  -n ecommerce -f ./helm-charts/microservices/values-prod.yaml
```


# Cleanup

```
microk8s kubectl delete namespace ecommerce 

microk8s helm3 uninstall microservices -n ecommerce 2>/dev/null

microk8s kubectl delete svc --all -A

microk8s kubectl delete secret -A -l owner=helm

microk8s kubectl delete pvc --all -A

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
