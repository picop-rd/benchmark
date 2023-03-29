#!/bin/bash -ex
PROXY_CONTROLLER_URL=$1

# proxy-http
curl ${PROXY_CONTROLLER_URL}/proxy/service-http/register -XPUT -H "Content-Type: application/json" -d "{\"endpoint\":\"http://proxy-http.service.svc.cluster.local:9000\"}"

# proxy-picop
curl ${PROXY_CONTROLLER_URL}/proxy/service-picop/register -XPUT -H "Content-Type: application/json" -d "{\"endpoint\":\"http://proxy-picop.service.svc.cluster.local:9000\"}"

