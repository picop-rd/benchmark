#!/bin/bash -e

while true; do
	docker stats --no-stream | awk '{ "date +\"%Y/%m/%d %T\"" | getline var; print var " " $0 }'
	sleep 1
done
