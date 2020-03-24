#!/bin/sh

set -e
# set -x

echo "Starting ETCD Hashring Simple Server..."
etcd-hashring-simple-server \
     --log-level="$LOG_LEVEL" \
     --etcd-url="$ETCD_URL" \
     --etcd-service-path="$ETCD_SERVICE_PATH" \
     --etcd-timeout-seconds="$ETCD_TIMEOUT_SECONDS"

