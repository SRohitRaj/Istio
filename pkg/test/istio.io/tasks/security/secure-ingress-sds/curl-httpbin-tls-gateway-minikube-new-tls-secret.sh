#!/bin/bash

INGRESS_HOST=$(kubectl -n istio-system get pod -l istio=ingressgateway -o jsonpath='{.items[0].status.hostIP}')
echo "INGRESS_HOST ${INGRESS_HOST}"
export INGRESS_HOST

SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.port==443)].nodePort}')
echo "SECURE_INGRESS_PORT ${SECURE_INGRESS_PORT}"
export SECURE_INGRESS_PORT

echo "curl -v -HHost:httpbin.example.com \
--resolve httpbin.example.com:${SECURE_INGRESS_PORT}:${INGRESS_HOST} \
--cacert httpbin.new.example.com/2_intermediate/certs/ca-chain.cert.pem \
https://httpbin.example.com:${SECURE_INGRESS_PORT}/status/418"

curl -v -HHost:httpbin.example.com \
--resolve httpbin.example.com:"$SECURE_INGRESS_PORT":"$INGRESS_HOST" \
--cacert httpbin.new.example.com/2_intermediate/certs/ca-chain.cert.pem \
https://httpbin.example.com:"$SECURE_INGRESS_PORT"/status/418