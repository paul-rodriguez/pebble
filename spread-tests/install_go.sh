#! /usr/bin/env bash

# This script installs the version of Go required by the go.mod file, as a
# snap.

set -e

GO_MOD="${SPREAD_PATH}/go.mod"
VERSION_PATTERN="[[:digit:]]+\.[[:digit:]]+"

GO_VERSION=$(cat "${GO_MOD}" \
    | grep -Ex "go\s+${VERSION_PATTERN}" \
    | grep -E --only-matching "${VERSION_PATTERN}")

snap wait system seed.loaded
snap install go --classic --channel=${GO_VERSION}
