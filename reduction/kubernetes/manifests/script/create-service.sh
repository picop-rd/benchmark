#!/bin/bash -e

startNum=$1
endNum=$2
manifests=""

for i in `seq $startNum $endNum`; do
	numStr=$(printf "%03d" $i)
	manifests=$manifests$'\n'$(cat ./template/service.yaml | sed -e "s/\$(ENV)/$numStr/g")
done

echo "$manifests"
