# Resource Reduction
You can measure the number of containers, CPU, Memory, and response time under load and autoscaled, and compare resource consumption with and without shared microservices using PiCoP.

# Preparation
[Common Steps](../docs/common.md)

```bash
cd kubernetes/manifests
kubectl apply -f namespace.yaml
kubectl apply -f proxy-http.yaml
./script/create-service.sh 1 <the number of environments> | kubectl -n service apply -f -
```

# Measurements
```bash
cd script
go run timertt/main.go <options>
```

## share: yes
```bash
go run timertt/main.go --url http://10.229.71.125:31000 --env-id main --req-per-sec 1000 --duration 300 --payload 1000 --client-num <the number of environments> --proxy --picop
```

## share: no
```bash
go run timertt/main.go --url http://10.229.71.125:31000 --env-id main --req-per-sec 1000 --duration 300 --payload 1000 --client-num 1
```

# Outputs
```bash
./benchmark.sh
```
