#!/bin/bash -eux

if [ "$#" -ne 3 ]; then
	echo "Usage: $0 <prefix> <share> <envs>"
	exit 1
fi

PREFIX=$1
SHARE=$2 # yes or no
ENVS=$3 # 10

RPS=1000
DURATION=300
INTERVAL=10

NAME="$PREFIX/$SHARE-$ENVS-$RPS-$DURATION-$INTERVAL"
TIMESTAMP=$(date +%s)

echo "NAME: $NAME"
echo "TIMESTAMP: $TIMESTAMP"

kubectl get deploy -n service
kubectl get deploy -n service -o json | jq '.items[].status.availableReplicas' | awk '{s += $1} END {print s}'

go run ./script/collect/main.go -name $NAME -timestamp $TIMESTAMP -dir ./data/resource/input -interval $INTERVAL -duration $DURATION kubectl top pod -n service --containers
go run ./script/parse/main.go -name $NAME -timestamp $TIMESTAMP -input ./data/resource/input -output ./data/resource/output --containers

NUM_PROXY=$(kubectl get deploy -n service --no-headers | grep proxy | awk '{ s += $4 } END { print s }')
NUM_SERVICE=$(kubectl get deploy -n service --no-headers | grep service | awk '{ s += $4 } END { print s }')

echo "proxy,service" >> ./data/resource/output/$NAME/$TIMESTAMP/org/container.csv
echo "$NUM_PROXY,$NUM_SERVICE" >> ./data/resource/output/$NAME/$TIMESTAMP/org/container.csv
