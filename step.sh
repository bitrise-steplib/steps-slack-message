#!/bin/bash
set -e
THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

go_package_name="github.com/bitrise-io/steps-slack-message"
full_package_path="${THIS_SCRIPT_DIR}/go/${go_package_name}"
mkdir -p "${full_package_path}"

rsync -avh --quiet "${THIS_SCRIPT_DIR}/" "${full_package_path}/"

export GOPATH="${THIS_SCRIPT_DIR}/go"
export GO15VENDOREXPERIMENT=1
go run "${full_package_path}/main.go"
