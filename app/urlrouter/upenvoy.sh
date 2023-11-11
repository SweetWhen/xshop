#!/bin/bash

docker run --rm --name xshop-urlrouter  -p 10000:10000 -v $(pwd)/envoy.yaml:/etc/envoy/envoy.yaml -v $(pwd)/share:/share envoyproxy/envoy:v1.21.4
