#!/bin/bash -e

pids=""
# container idの取得
for container_id in $(docker ps | awk '!/CONTAINER/ { print $1 }'); do
	echo $container_id

	# pidの取得
	pid=$(docker inspect $container_id | awk '/"Pid"/ { print $2 }')
	pids=$pids$pid
done
pids=${pids/%?/}
echo $pids

pidstat 1 -udrl -p $pids
