#!/bin/bash -e

protocol=$1
environment=$2
port=$3

cat ./template/service.yaml | sed -e "s/\$(PROTO)/$protocol/g" | sed -e "s/\$(ENV)/$environment/g" | sed -e "s/\$(PORT)/$port/g"
