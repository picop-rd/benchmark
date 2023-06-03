# Proxy Communication Delay
You can compare HTTP request response time with/without PiCoP/Istio proxy.

# Preparation
[Common Steps](../docs/common.md)

```
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
```
cd script
go run timertt/main.go <options>
```

Please change --client-num from 1~64
### base
```
go run timertt/main.go --url http://service-istio.service-istio.svc.cluster.local:32001 --prefix base --env-id main --req-per-sec 1000 --duration 10 --client-num 1 --payload 1000
```
### base+proposed
```
go run timertt/main.go --url http://service-istio.service-istio.svc.cluster.local:31002 --prefix base+proposed --env-id main --req-per-sec 1000 --duration 10 --client-num 1 --payload 1000 --picop
```
### base+gw+istio
```
go run timertt/main.go --url http://service-istio.service-istio.svc.cluster.local:30001 --prefix base+gw+istio --env-id main --req-per-sec 1000 --duration 10 --client-num 1 --payload 1000
```
### base+gw+proposed
```
go run timertt/main.go --url http://proxy-both.service.svc.cluster.local:30002 --prefix base+gw+proposed --env-id main --req-per-sec 1000 --duration 10 --client-num 1 --payload 1000 --picop
```

# Outputs
CSV format
```
count: Request Index
latency: Response Time[ns]
start: Request Sending Time[ns] (Accuracy in ms)
end: Response Receiving Time[ns] (Accuracy in ms)
```

