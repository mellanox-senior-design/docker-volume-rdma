#!/bin/bash

function red() {
	echo $'\033[1;31m'"$@"$'\033[m'
}

function green() {
	echo $'\033[1;32m'"$@"$'\033[m'
}

GIT_BRANCH=$(
	git rev-parse --abbrev-ref HEAD |
	tr '[:upper:]' '[:lower:]' |
	sed 's#[^a-z0-9._-]#-#'
)

set -e

green Unit Tests
if ! docker build --tag docker-volume-rdma:"$GIT_BRANCH" .; then
	red 'Failed to run unit tests'
fi
docker rmi docker-volume-rdma:"$GIT_BRANCH"

green Benchmarks
echo # make newline space
for scenario in $(cd ./benchmarking; ls -d */ | grep -v '^scenarios/?$' | sed 's#/##'); do
	echo # make space
	green "Running Scenario: $scenario"
	echo # make space
	sleep 5
	if ! ./benchmarking/test.sh "$scenario"; then
		red 'Failed to run benchmarks'
	fi
done
