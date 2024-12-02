#!/bin/bash -eux

if [ "$#" -ne 2 ]; then
	echo "Usage: $0 <share> <prefix>"
	exit 1
fi

# SHARE=yes
# SHARE=no
SHARE=$1
PREFIX=$2

ENVS=100

URL="http://10.229.71.125:31000"
DURATION=300
INTERVAL=10
RPS=1000

NAME="$PREFIX/$SHARE-$ENVS-$RPS-$DURATION-$INTERVAL"
TIMESTAMP=$(date +%s)

echo "NAME: $NAME"
echo "TIMESTAMP: $TIMESTAMP"

kubectl get deploy -n service
kubectl get deploy -n service -o json | jq '.items[].status.availableReplicas' | awk '{s += $1} END {print s}'

go run ./script/collect/main.go -name $NAME -timestamp $TIMESTAMP -dir ./data/resource/input -interval $INTERVAL -duration $DURATION kubectl top pod -n service --containers
go run ./script/parse/main.go -name $NAME -timestamp $TIMESTAMP -input ./data/resource/input -output ./data/resource/output --containers
