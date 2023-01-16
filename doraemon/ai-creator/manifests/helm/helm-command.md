### namespace
```shell
kubectl create namespace doraemon
```

### helm Chart repo
```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
```

### helm install redis Chart
```shell
# https://github.com/bitnami/charts/tree/main/bitnami/redis
helm install --namespace doraemon ai-creator bitnami/redis
# use https://github.com/bitnami/charts/blob/main/bitnami/redis/values.yaml for defined env params
kubectl apply -f manifests/app-secret.yaml
helm install --namespace doraemon ai-creator bitnami/redis --values manifests/helm/bitnami-redis-values.yml
```

### helm install redis-cluster Chart
```shell
# https://github.com/bitnami/charts/tree/main/bitnami/redis-cluster
helm install --namespace doraemon ai-creator bitnami/redis-cluster
# use https://github.com/bitnami/charts/blob/main/bitnami/redis-cluster/values.yaml for defined env params 
kubectl apply -f manifests/app-secret.yaml
helm install --namespace doraemon ai-creator bitnami/redis-cluster --values manifests/helm/bitnami-redis-cluster-values.yml
```

### helm uninstall/delete Chart
```shell
helm delete --namespace doraemon ai-creator
```

### k8s kind
```shell
kubectl apply -f manifests/app-configmap.yaml 
#kubectl delete -f manifests/app-configmap.yaml 
kubectl apply -f manifests/app-deployment.yaml
#kubectl delete -f manifests/app-deployment.yaml

kubectl get all --namespace doraemon
kubectl get pods -o wide -w --namespace doraemon
kubectl get svc -w -o wide --namespace doraemon

#kubectl logs -n doraemon -f -c ai-creator-container ai-creator-deployment-77d9d85dd-wl4rx
#kubectl logs -n doraemon -f -c ai-creator-container ai-creator-deployment-77d9d85dd-rxgjf
# SLS or E(C)LK
```


### minikube service
```shell
minikube service list
minikube service -n doraemon ai-creator-service
```

### scaling
```
kubectl scale -n doraemon --replicas 2 deployment/ai-creator-deployment
```


### balancing

### reference
1. https://cloud.google.com/blog/products/containers-kubernetes/your-guide-kubernetes-best-practices
2. https://learnk8s.io/kubernetes-long-lived-connections