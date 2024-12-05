# Proxy Communication Delay
You can compare HTTP request response time with/without PiCoP/Istio proxy.

# Preparation
[Common Steps](../docs/common.md)

```bash
cd script
./register-proxies.sh <Proxy Controller URL>
./register-routes.sh <Proxy Controller URL>

cd kubernetes/manifests
kubectl apply -f namespace.yaml 
kubectl apply -f proxy-http.yaml
kubectl apply -f proxy-both.yaml
./script/create-service.sh http main | kubectl -n service apply -f -
./script/create-service.sh http feature-1 | kubectl -n service apply -f -
./script/create-service.sh http main | kubectl -n service-istio apply -f -
./script/create-service.sh http feature-1 | kubectl -n service-istio apply -f -

# Require istioctl
# See https://istio.io/latest/docs/setup/getting-started/
cd istio
istioctl install --set profile=default -y
kubectl delete -f save/istio-ingressgateway-hpa.yaml
kubectl delete -f save/istiod-hpa.yaml
# Please change replicas of Deployments `istio-ingressgateway` and `istiod` to 5 (max value of hpa).

cd istio
kubectl apply -f namespace.yaml
kubectl apply -f gateway.yaml
kubectl apply -f vs.yaml
kubectl apply -f ingressgateway-svc.yaml
# Please append etc/hosts to that of the benchmark machine.(Please replace the IP address 192.168.0.2 with that of the Kubernetes Cluster.)
```

# Measurements
Please change --client-num from 1~64
### base
```bash
./benchmark.sh base <PREFIX> <CLIENT>
```
### base+picop
```bash
./benchmark.sh base+picop <PREFIX> <CLIENT>
```
### base+gw+istio
```bash
./benchmark.sh base+gw+istio <PREFIX> <CLIENT>
```
### base+gw+picop
```bash
./benchmark.sh base+gw+picop <PREFIX> <CLIENT>
```

# Outputs
CSV format
```bash
count: Request Index
latency: Response Time[ns]
start: Request Sending Time[ns] (Accuracy in ms)
end: Response Receiving Time[ns] (Accuracy in ms)
```

