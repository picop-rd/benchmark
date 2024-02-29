#!/bin/bash -eux

if [ "$#" -ne 1 ]; then
	echo "Usage: $0 <prefix>"
	exit 1
fi

PREFIX=$1

URL="http://10.229.71.125:32001"
ENVID=main
DURATION=300
CONN=10
RPS=50
PAYLOAD=1000

INTERVAL=10

NAME="$PREFIX/$ENVID-$CONN-$RPS-$DURATION-$PAYLOAD-$INTERVAL"
TIMESTAMP=$(date +%s)

CMD="go run ./script/main.go --url $URL --env-id $ENVID --req-per-sec $RPS --duration $DURATION --client-num $CONN --payload $PAYLOAD --picop"

echo "NAME: $NAME"
echo "TIMESTAMP: $TIMESTAMP"

cleanup() {
    echo "SIGINT received, cleaning up..."
    ssh onoe-benchmark "pkill -2 --echo 'go'"
    ssh onoe-benchmark "pkill -2 --echo 'tee'"
    ssh onoe-benchmark "pkill -2 --echo 'exec-cmd.sh'"

    kill $(jobs -p)
    exit 1
}

trap 'cleanup' SIGINT

ssh onoe-benchmark "cd benchmark/consumption && ./exec-cmd.sh '$CMD' ./data/$NAME/cmd $TIMESTAMP" &

go run ./collect/main.go -name $NAME -timestamp $TIMESTAMP -dir ./data/input -interval $INTERVAL -duration $DURATION kubectl top pod -n service-proxy &

wait

mkdir -p ./data/$NAME/cmd
scp onoe-benchmark:benchmark/consumption/data/$NAME/cmd/$TIMESTAMP.txt ./data/$NAME/cmd/$TIMESTAMP.txt
go run ./parse/main.go -name $NAME -timestamp $TIMESTAMP -input ./data/input -output ./data/output
