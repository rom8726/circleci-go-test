#!/usr/bin/env bash

sleep 10;
curl -X POST -u Administrator:password -d 'name=test' -d 'ramQuotaMB=100' -d 'authType=none' -d 'proxyPort=11216' http://127.0.0.1:8091/pools/default/buckets