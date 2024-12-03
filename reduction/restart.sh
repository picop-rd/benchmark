#!/bin/bash -ux

ENVS=$1

cd kubernetes/manifests

kubectl delete -f namespace.yaml
kubectl apply -f namespace.yaml
kubectl apply -f proxy-http.yaml
./script/create-service.sh 1 $ENVS | kubectl -n service apply -f -
