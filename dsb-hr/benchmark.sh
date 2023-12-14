#!/bin/bash -eux

URL="http://10.229.71.125:31000"
LUA="./scripts/hotel-reservation/mixed-workload_type_1.lua"
DURATION=300
CONN=10
RPS=1000
CMD="./wrks/wrk -D exp -t 1 -c $CONN -d $DURATION -L -R $RPS -S -r -s $LUA $URL"

INTERVAL=10

NAME=mixed-$CONN-$RPS-$DURATION-$INTERVAL
TIMESTAMP=$(date +%s)

ssh onoe-benchmark "cd benchmark/dsb-hr && ./latency/collect.sh '$CMD' ./latency/data/$NAME $TIMESTAMP" &

go run ./resource/collect/main.go -name $NAME -timestamp $TIMESTAMP -dir ./resource/data/input -interval $INTERVAL -duration $DURATION kubectl top pod -n dsb-hr &

wait

mkdir -p ./latency/data/$NAME
scp onoe-benchmark:benchmark/dsb-hr/latency/data/$NAME/$TIMESTAMP.txt ./latency/data/$NAME/$TIMESTAMP.txt
go run ./resource/parse/main.go -name $NAME -timestamp $TIMESTAMP -input ./resource/data/input -output ./resource/data/output
