# benchmark
Benchmark

## 計測
- [Proxy Communication Delay](./latency/README.md)
- [Resource Reduction](./manyenvs/README.md)
## 各計測で共通の準備
```
cd ../demo/kubernetes/manifests
kubectl apply -f namespace.yaml proxy-controller.yaml
```
