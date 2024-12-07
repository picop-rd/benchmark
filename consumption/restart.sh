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
kubectl apply -n service -f config.yaml
kubectl apply -f service.yaml
kubectl apply -n service-istio -f config.yaml
kubectl apply -f service-istio.yaml
