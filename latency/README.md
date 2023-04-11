PiCoPやIstioのプロキシのあるなしでHTTPリクエストのレスポンスタイムを比較する

# 準備
[共通の準備](../README.md)
```
cd kubernetes/manifests
kubectl apply -f namespace.yaml proxy-http.yaml proxy-picop.yaml
./script/create-service.sh http main | kubectl -n service apply -f -
./script/create-service.sh picop main | kubectl -n service apply -f -

# Require istioctl
# See https://istio.io/latest/docs/setup/getting-started/
istioctl install --set profile=default -y
kubectl delete -f save/istio-ingressgateqay-hpa.yaml istiod-hpa.yaml

cd istio
kubectl apply namespave.yaml gateway.yaml vs.yaml

cd scriot
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

