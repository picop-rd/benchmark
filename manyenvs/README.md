# Resource Reduction
You can measure the number of containers under load and autoscaled, and compare resource consumption with and without shared microservices using PiCoP.

# Preparation
[Common Steps](../docs/common.md)

```
cd kubernetes/manifests
kubectl apply -f namespace.yaml
kubectl apply -f proxy-http.yaml
./script/create-service.sh 1 <the number of environments> | kubectl -n service apply -f -
```

# Measurements
```
cd script
go run timertt/main.go <options>
```

## share: yes
```
go run timertt/main.go --url http://<Kubernetes Cluster IP address>:30100 --env-id main --req-per-sec 1000 --duration 300 --client-num <the number of environments> --payload 1000 --proxy --picop
```

## share: no
```
go run timertt/main.go --url http://<Kubernetes Cluster IP address>:32 --env-id main --req-per-sec 1000 --duration 300 --client-num 1 --payload 1000
```

# Outputs
You can use `kubectl` and `k9s` to measure the number of containers.
