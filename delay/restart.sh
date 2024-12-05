#!/bin/bash -ux

cd kubernetes/manifests

kubectl delete -f namespace.yaml
kubectl apply -f namespace.yaml
kubectl apply -f proxy-http.yaml
kubectl apply -f proxy-both.yaml
./script/create-service.sh http main | kubectl -n service apply -f -
./script/create-service.sh http feature-1 | kubectl -n service apply -f -
./script/create-service.sh http main | kubectl -n service-istio apply -f -
./script/create-service.sh http feature-1 | kubectl -n service-istio apply -f -

cd ../../istio
kubectl delete -f namespace.yaml
kubectl apply -f namespace.yaml
kubectl apply -f gateway.yaml
kubectl apply -f vs.yaml
kubectl apply -f ingressgateway-svc.yaml
