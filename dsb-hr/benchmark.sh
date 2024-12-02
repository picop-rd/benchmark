#!/bin/bash -eux

if [ "$#" -ne 1 ]; then
	echo "Usage: $0 <prefix>"
	exit 1
fi

PREFIX=$1

# LUA_NAME=mixed-workload_type_1
# LUA_NAME=recommend
LUA_NAME=reserve
# LUA_NAME=search_hotel
# LUA_NAME=user_login

# base or picop
TYPE=$(cd ../../DeathStarBench && git branch --contains | cut -d " " -f 2)

URL="http://10.229.71.125:31000"
LUA_DIR="./scripts/hotel-reservation/$LUA_NAME.lua"
DURATION=300
CONN=10
RPS=40
CMD="./wrks/wrk -D exp -t 1 -c $CONN -d $DURATION -L -R $RPS -S -r -s $LUA_DIR $URL"

INTERVAL=10

NAME="$PREFIX/$LUA_NAME-$CONN-$RPS-$DURATION-$INTERVAL"
TIMESTAMP=$(date +%s)

echo "TYPE: $TYPE"
echo "NAME: $NAME"
echo "TIMESTAMP: $TIMESTAMP"

cleanup() {
    echo "SIGINT received, cleaning up..."
    ssh onoe-benchmark-1 "pkill -2 --echo 'wrk'"
    ssh onoe-benchmark-1 "pkill -2 --echo 'tee'"
    ssh onoe-benchmark-1 "pkill -2 --echo 'collect.sh'"

    kill $(jobs -p)
    exit 1
}

trap 'cleanup' SIGINT

ssh onoe-benchmark-1 "cd benchmark/dsb-hr && ./latency/collect.sh '$CMD' ./latency/data/$TYPE/$NAME $TIMESTAMP" &

go run ./resource/collect/main.go -name $NAME -timestamp $TIMESTAMP -dir ./resource/data/input/$TYPE -interval $INTERVAL -duration $DURATION kubectl top pod -n dsb-hr --containers &

wait

mkdir -p ./latency/data/$TYPE/$NAME
scp onoe-benchmark-1:benchmark/dsb-hr/latency/data/$TYPE/$NAME/$TIMESTAMP.txt ./latency/data/$TYPE/$NAME/$TIMESTAMP.txt
go run ./resource/parse/main.go -name $NAME -timestamp $TIMESTAMP -input ./resource/data/input/$TYPE -output ./resource/data/output/$TYPE --containers
