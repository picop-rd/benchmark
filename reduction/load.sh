#!/bin/bash -eux

if [ "$#" -ne 2 ]; then
	echo "Usage: $0 <share> <envs>"
	exit 1
fi

SHARE=$1 # yes or no
ENVS=$2 # 10
DURATION=600
RPS=1000

URL=http://10.229.71.125
if [ "$SHARE" = "yes" ]; then
	URL="$URL:30100"
else
	URL="$URL:32"
fi

CMD="/usr/local/go/bin/go run timertt/main.go --url $URL --env-id main --req-per-sec $RPS --duration $DURATION --payload 1000"

cleanup() {
    echo "SIGINT received, cleaning up..."
    ssh onoe-benchmark-1 "pkill -2 --echo 'go'"

    kill $(jobs -p)
    exit 1
}

trap 'cleanup' SIGINT

if [ "$SHARE" = "yes" ]; then
	CMD="$CMD --client-num $ENVS --proxy --picop"
else
	CMD="$CMD --client-num $ENVS"
fi

ssh onoe-benchmark-1 "cd benchmark/reduction/script && $CMD"
