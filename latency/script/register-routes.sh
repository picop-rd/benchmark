#!/bin/bash -ex
PROXY_CONTROLLER_URL=$1

curl ${PROXY_CONTROLLER_URL}/routes -XPUT -H "Content-Type: application/json" -d "
[
	{\"proxy_id\":\"service-http\",     \"env_id\": \"feature-1\", \"destination\": \"service-http-feature-1.service.svc.cluster.local:80\"},
	{\"proxy_id\":\"service-picop\",     \"env_id\": \"feature-1\", \"destination\": \"service-picop-feature-1.service.svc.cluster.local:80\"}
]
"
