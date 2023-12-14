#!/bin/bash -eux

LUA_NAME=mixed-workload_type_1
# LUA_NAME=recommend
# LUA_NAME=reserve
# LUA_NAME=search_hotel
# LUA_NAME=user_login

URL="http://10.229.71.125:31000"
LUA_DIR="./scripts/hotel-reservation/$LUA_NAME.lua"
DURATION=300
CONN=10
RPS=1000
CMD="./wrks/wrk -D exp -t 1 -c $CONN -d $DURATION -L -R $RPS -S -r -s $LUA_DIR $URL"

INTERVAL=10

NAME=$LUA_NAME-$CONN-$RPS-$DURATION-$INTERVAL
TIMESTAMP=$(date +%s)

cleanup() {
    echo "SIGINT received, cleaning up..."
    ssh onoe-benchmark "pkill -2 --echo 'wrk'"
    ssh onoe-benchmark "pkill -2 --echo 'tee'"
    ssh onoe-benchmark "pkill -2 --echo 'collect.sh'"

    kill $(jobs -p)
    exit 1
}

trap 'cleanup' SIGINT

ssh onoe-benchmark "cd benchmark/dsb-hr && ./latency/collect.sh '$CMD' ./latency/data/$NAME $TIMESTAMP" &

go run ./resource/collect/main.go -name $NAME -timestamp $TIMESTAMP -dir ./resource/data/input -interval $INTERVAL -duration $DURATION kubectl top pod -n dsb-hr &

wait

mkdir -p ./latency/data/$NAME
scp onoe-benchmark:benchmark/dsb-hr/latency/data/$NAME/$TIMESTAMP.txt ./latency/data/$NAME/$TIMESTAMP.txt
go run ./resource/parse/main.go -name $NAME -timestamp $TIMESTAMP -input ./resource/data/input -output ./resource/data/output
