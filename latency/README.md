PiCoPやIstioのプロキシのあるなしでHTTPリクエストのレスポンスタイムを比較する

# 準備
[共通の準備](../README.md)
```
cd kubernetes/manifests
kubectl apply -f namespace.yaml 
kubectl apply -f proxy-http.yaml
kubectl apply -f proxy-picop.yaml
./script/create-service.sh http main | kubectl -n service apply -f -
./script/create-service.sh picop main | kubectl -n service apply -f -

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

cd script
./register-proxies.sh <Proxy Controller URL>
./register-routes.sh <Proxy Controller URL>
```

# 計測

```
cd script
go run timertt/main.go <条件に基づいて設定したオプション>
```

# 出力
CSV形式で出力される。各カラムは
```
count: 番号
latency: リクエストを送信してからレスポンスが帰ってくるまでの時間(ns)
start: リクエスト送信時間(ns)
end: レスポンス受信時間(ns)
```
