#! /usr/bin/env bash

# Installs gcc as a deb
#
# This is required because gcc is one of the dependencies of the test
# framework.

set -e

apt-get update
apt-get install gcc -y

