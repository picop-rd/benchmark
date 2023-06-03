# Resource Reduction
You can measure proxy resource consumption (vCPU and memory).

# Preparation
[Common Steps](../docs/common.md)

```bash
cd kubernetes/manifests
kubectl apply -f namespace.yaml
kubectl apply -f proxy-http.yaml
kubectl apply -f service.yaml
```

# Measurements
Please change --req-per-sec from 1~10000
```bash
cd script
go run --url http://<Kubernetes Cluster URL>:32001 --env-id main --req-per-sec 1000 --duration 90 --client-num 1 --payload 1000 --picop
```

# Outputs
You can use `kubectl top` and `k9s` to measure the number of vCPUs and memory.
