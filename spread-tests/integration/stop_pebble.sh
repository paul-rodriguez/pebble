#! /usr/bin/env bash

set -e

declare PEBBLE_PID=$(pgrep -x pebble)

if [ -n "$PEBBLE_PID" ]; then
    echo "Killing pebble"
    kill -s TERM $PEBBLE_PID
fi

rm -f \
    "${TASK_DIR}/.pebble.socket" \
    "${TASK_DIR}/.pebble.state" \
    "${TASK_DIR}/.pebble.socket.untrusted" \
    "${TASK_DIR}/dummy_service_started.tok"
