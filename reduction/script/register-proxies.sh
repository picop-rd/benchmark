#!/bin/bash -ex
PROXY_CONTROLLER_URL=$1

# proxy-http
curl ${PROXY_CONTROLLER_URL}/proxy/service/register -XPUT -H "Content-Type: application/json" -d "{\"endpoint\":\"http://proxy-http.service.svc.cluster.local:9000\"}"
