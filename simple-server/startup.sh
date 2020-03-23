#!/bin/sh

set -e
# set -x

echo "Starting ETCD Hashring Simple Server..."
etcd-hashring-simple-server \
     --log-level=$LOG_LEVEL \
     --

