#! /usr/bin/env bash

set -e

pushd ${SPREAD_PATH}
    su ubuntu -s /bin/sh -c "go test $1"
popd

