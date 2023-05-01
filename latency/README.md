PiCoPやIstioのプロキシのあるなしでHTTPリクエストのレスポンスタイムを比較する

# 準備
[共通の準備](../README.md)
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
istioctl install --set profile=default -y
kubectl delete -f save/istio-ingressgateway-hpa.yaml
kubectl delete -f save/istiod-hpa.yaml
# Deploymentのistio-ingressgateway, istiodのreplicasを5へ変更(hpaのmax値)

cd istio
kubectl apply -f namespace.yaml
kubectl apply -f gateway.yaml
kubectl apply -f vs.yaml
kubectl apply -f ingressgateway-svc.yaml
# etc/hostsの内容をベンチマーク実行マシンで適用する
```

# 計測

```
cd script
go run timertt/main.go <条件に基づいて設定したオプション>
```

- --client-numを1~64で変更して計測
- URLに含まれるport番号は各serviceのNodePortを参照して設定
### base
```
go run timertt/main.go --url http://service-istio.service-istio.svc.cluster.local:30830 --prefix no-http-main-2-1Gi-100m-128Mi --env-id main --req-per-sec 1000 --duration 10 --client-num 1 --payload 1000
```
### base+proposed
```
go run timertt/main.go --url http://service-istio.service-istio.svc.cluster.local:30100 --prefix proxy-http-main-2-1Gi-100m-128Mi --env-id main --req-per-sec 1000 --duration 10 --client-num 1 --payload 1000 --picop
```
### base+gw+istio
```
go run timertt/main.go --url http://service-istio.service-istio.svc.cluster.local:30031 --prefix istio-http-main-2-1Gi-100m-128Mi --env-id main --req-per-sec 1000 --duration 10 --client-num 1 --payload 1000
```
### base+gw+proposed
```
go run timertt/main.go --url http://proxy-both.service.svc.cluster.local:30603 --prefix both-http-main-2-1Gi-100m-128Mi --env-id main --req-per-sec 1000 --duration 10 --client-num 1 --payload 1000 --picop
```

# 出力
CSV形式で出力される。各カラムは
```
count: 番号
latency: リクエストを送信してからレスポンスが帰ってくるまでの時間(ns)
start: リクエスト送信時間(ns)
end: レスポンス受信時間(ns)
```

