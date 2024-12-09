#!/bin/bash -eux

if [ "$#" -ne 3 ]; then
	echo "Usage: $0 <type> <prefix> <rps>"
	exit 1
fi

TYPE=$1
PREFIX=$2
RPS=$3

# TYPE=base
# TYPE=base+picop
# TYPE=base+gw+istio
# TYPE=base+gw+picop

DURATION=300
CLIENT=1
INTERVAL=10

TIMESTAMP=$(date +%s)
NAME="$PREFIX/$TYPE/$CLIENT-$RPS-$DURATION/$TIMESTAMP"

if [ "$TYPE" = "base" ]; then
    echo "Invalid type: $TYPE"
    exit 1
    # URL="http://service-istio.service-istio.svc.cluster.local:32001"
    # OPTION=""
    # NS="service"
elif [ "$TYPE" = "base+picop" ]; then
    echo "Invalid type: $TYPE"
    exit 1
    # URL="http://service-istio.service-istio.svc.cluster.local:31002"
    # OPTION="--picop"
    # NS="service"
elif [ "$TYPE" = "base+gw+istio" ]; then
    URL="http://service-istio.service-istio.svc.cluster.local:30001"
    OPTION=""
    NS="service-istio"
elif [ "$TYPE" = "base+gw+picop" ]; then
    URL="http://proxy-both.service.svc.cluster.local:30002"
    OPTION="--picop"
    NS="service"
else
    echo "Invalid type: $TYPE"
    exit 1
fi

CMD="/usr/local/go/bin/go run script/main.go --url $URL --env-id main --req-per-sec $RPS --duration $DURATION --client-num $CLIENT --payload 1000 $OPTION"

cleanup() {
    echo "SIGINT received, cleaning up..."
    ssh onoe-benchmark-1 "pkill -2 --echo 'go'"
    # ssh onoe-benchmark-2 "pkill -2 --echo 'go'"

    kill $(jobs -p)
    exit 1
}

trap 'cleanup' SIGINT

if [ "$RPS" -gt 0 ]; then
    ssh onoe-benchmark-1 "cd benchmark/consumption && $CMD" &
    # ssh onoe-benchmark-2 "cd benchmark/consumption && $CMD" &
fi

go run ./collect/main.go -name ./data/input/$NAME -timestamp $TIMESTAMP -interval $INTERVAL -duration $DURATION kubectl top pod -n $NS --containers &

wait

go run ./parse/main.go -name $NAME -timestamp $TIMESTAMP -input ./data/input -output ./data/output

echo "NAME: $NAME"
echo "CMD: $CMD"
echo "TIMESTAMP: $TIMESTAMP"
