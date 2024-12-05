#!/bin/bash -eux

if [ "$#" -ne 1 ]; then
	echo "Usage: $0 <type> <prefix> <client>"
	exit 1
fi

PREFIX=$1
CLIENT=$2

# TYPE=base
# TYPE=base+picop
# TYPE=base+gw+istio
# TYPE=base+gw+picop

DURATION=10
RPS=1000

TIMESTAMP=$(date +%s)
NAME="data/$PREFIX/$TYPE/$CLIENT-$RPS-$DURATION/$TIMESTAMP"

if [ "$TYPE" = "base" ]; then
	URL="http://service-istio.service-istio.svc.cluster.local:32001"
	OPTION=""
elif [ "$TYPE" = "base+picop" ]; then
	URL="http://service-istio.service-istio.svc.cluster.local:31002"
	OPTION="--picop"
elif [ "$TYPE" = "base+gw+istio" ]; then
	URL="http://service-istio.service-istio.svc.cluster.local:30001"
	OPTION=""
elif [ "$TYPE" = "base+gw+picop" ]; then
	URL="hhttp://proxy-both.service.svc.cluster.local:30002"
	OPTION="--picop"
else
	echo "Invalid type: $TYPE"
	exit 1
fi

CMD="/usr/local/go/bin/go run timertt/main.go --url $URL --prefix $NAME --env-id main --req-per-sec $RPS --duration $DURATION --client-num $CLIENT --payload 1000 $OPTION"

echo "NAME: $NAME"
echo "CMD: $CMD"
echo "TIMESTAMP: $TIMESTAMP"

cleanup() {
    echo "SIGINT received, cleaning up..."
    ssh onoe-benchmark-1 "pkill -2 --echo 'go'"

    kill $(jobs -p)
    exit 1
}

trap 'cleanup' SIGINT

ssh onoe-benchmark-1 "cd benchmark/consumption && $CMD"

mkdir -p ./data/$PREFIX/$TYPE/$CLIENT-$RPS-$DURATION
scp onoe-benchmark-1:benchmark/consumption/$NAME.csv ./$NAME.csv
