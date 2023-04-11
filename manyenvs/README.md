各条件でコンテナに負荷をかけてオートスケールした後のコンテナ数を計測する

# 準備
[共通の準備](../README.md)
```
cd kubernetes/manifests
kubectl apply -f namespace.yaml proxy-http.yaml
./script/create-service.sh 1 <環境数> | kubectl -n service apply -f -
```

# 計測

```
cd script
go run timertt/main.go <条件に基づいて設定したオプション>
```

# 出力
`kubectl`や`k9s`を利用して目視でコンテナ数を記録する
