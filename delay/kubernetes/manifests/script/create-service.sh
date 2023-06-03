#!/bin/bash -e

protocol=$1
environment=$2

cat ./template/service.yaml | sed -e "s/\$(PROTO)/$protocol/g" | sed -e "s/\$(ENV)/$environment/g"
