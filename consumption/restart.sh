#!/bin/bash -ux

cd istio
kubectl delete -f namespace.yaml
kubectl apply -f namespace.yaml
kubectl apply -f gateway.yaml
kubectl apply -f vs.yaml
kubectl apply -f ingressgateway-svc.yaml

cd ../kubernetes/manifests

kubectl delete -f namespace.yaml
kubectl apply -f namespace.yaml
kubectl apply -f proxy-both.yaml
./script/create-service.sh http main 32001 | kubectl -n service apply -f -
./script/create-service.sh http main 32002 | kubectl -n service-istio apply -f -
