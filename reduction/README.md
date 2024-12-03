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
./load.sh yes <ENVS>
```

## share: no
```bash
./load.sh no <ENVS>
```

# Outputs
```bash
./benchmark.sh <PREFIX> <yes or no> <ENVS>
```
