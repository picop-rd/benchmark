#!/bin/bash -eux

SHARE=$1

URL=http://10.229.71.125:31000
ENVS=100

CMD="go run timertt/main.go --url $URL --env-id main --req-per-sec 1000 --duration 300 --payload 1000"

cleanup() {
    echo "SIGINT received, cleaning up..."
    ssh onoe-benchmark-1 "pkill -2 --echo 'go'"

    kill $(jobs -p)
    exit 1
}

trap 'cleanup' SIGINT

if [ "$SHARE" = "yes" ]; then
	CMD="$CMD --client-num $ENVS --proxy --picop"
	ssh onoe-benchmark-1 "cd benchmark/reduction/script && $CMD"
else
	CMD="$CMD --client-num 1"
	ssh onoe-benchmark-1 "cd benchmark/reduction/script && $CMD"
fi
