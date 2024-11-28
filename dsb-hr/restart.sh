#!/bin/bash -ux

# # dsb-hr namespaceのdeploymentのname一覧を取得
# deployments=$(kubectl get deployments -n dsb-hr -o jsonpath='{.items[*].metadata.name}')

# # 取得したdeploymentそれぞれでkubectl rollout restartコマンドを実行
# for deployment in $deployments
# do
#     kubectl rollout restart deployment/$deployment -n dsb-hr
#     echo "Restarted deployment: $deployment"
# done

kubectl delete -n dsb-hr -Rf ../../DeathStarBench/hotelReservation/kubernetes
kubectl apply -n dsb-hr -Rf ../../DeathStarBench/hotelReservation/kubernetes
