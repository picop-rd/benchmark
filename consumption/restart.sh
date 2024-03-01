#!/bin/bash -eux

# # dsb-hr namespaceのdeploymentのname一覧を取得
# deployments=$(kubectl get deployments -n dsb-hr -o jsonpath='{.items[*].metadata.name}')

# # 取得したdeploymentそれぞれでkubectl rollout restartコマンドを実行
# for deployment in $deployments
# do
#     kubectl rollout restart deployment/$deployment -n dsb-hr
#     echo "Restarted deployment: $deployment"
# done

kubectl delete -Rf ./kubernetes/manifests
kubectl apply -Rf ./kubernetes/manifests
