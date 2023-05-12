#! /usr/bin/env bash

# This script wraps the command and setup needed to run pebble from source

set -e

export PEBBLE="${TASK_DIR}"
SRC="${SPREAD_PATH}/cmd/pebble"

echo "Running command \"go run $SRC $@\""

go run $SRC $@

