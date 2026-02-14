
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
microk8s kubectl logs -f <pod-name> -n ecommerce


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
"./helm,"\
"./k8s,"\
"./tmp" \
"./main.go" \
"./tmp" \
"output.txt"

../../../tools/llmcode/export_structure.sh \
  ./ \
  tmp \
  output.txt


../tools/llmcode/import_code.sh input.txt ./

```