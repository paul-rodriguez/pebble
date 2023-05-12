#! /usr/bin/env bash

set -e

"${SUITE_DIR}/pebble_wrapper.sh" run &> /dev/null &
disown $!

# Spin on this file to ensure the dummy service has been started
until [ -e "${TASK_DIR}/dummy_service_started.tok" ]
do
    echo "Waiting for dummy_service_started.tok to exist"
    sleep 1
done

echo "Found dummy_service_started.tok, preparation finished"
sleep 1

