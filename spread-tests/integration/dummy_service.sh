#! /usr/bin/env bash

set -e

touch "${TASK_DIR}/dummy_service_started.tok"

while true
do
    echo "This is a dummy service log"
    sleep 2
done


