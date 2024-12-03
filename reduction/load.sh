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

CMD="/usr/local/go/bin/go run timertt/main.go --env-id main --req-per-sec $RPS --duration $DURATION --payload 1000 --env-total $ENVS"

cleanup() {
    echo "SIGINT received, cleaning up..."
    ssh onoe-benchmark-1 "pkill -2 --echo 'go'"

    kill $(jobs -p)
    exit 1
}

trap 'cleanup' SIGINT

if [ "$SHARE" = "yes" ]; then
	CMD="$CMD --picop"
fi

for NUM in $(seq 1 $ENVS); do
	if [ "$SHARE" = "yes" ]; then
		ADDR="$URL:30100"
	else
		TMP_NUM="000$NUM"
		ADDR="$URL:32${TMP_NUM: -3}"
	fi
	echo "ADDR: $ADDR"
	ssh onoe-benchmark-1 "cd benchmark/reduction/script && $CMD --env-num $NUM --url $ADDR" &
done

wait
