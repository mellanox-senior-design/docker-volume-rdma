#! /bin/bash

function testOut() {
    echo "[INFO] $1"
}

function testWarnOut() {
    echo "[WARN] $1"
}

function testErrorOut() {
    echo "[ERROR] $1"
}

function dcDown() {
    testOut "Stopping..."
    docker-compose -f docker-compose.yml -f "$1" down -v
}

function dcUp() {
    testOut "Bringing up Test Fixture [$1]..."
    if [ -f "$2" ]; then
        docker-compose -f docker-compose.yml -f "$2" up --abort-on-container-exit --force-recreate
    else
        testWarnOut "File, $2, does not exist. Skipping..."
    fi
}

function dcCollect() {
    testOut "Collecting..."
    if [[ -f /tmp/bench_results/result.json ]]; then
        mv /tmp/bench_results/result.json "$1"
        testOut "$(cat $1)"
    fi
}

function dcCollectAll() {
    testOut "Collecting All..."
    ./generate_report.py > /tmp/bench_results/result.json
    if [ "$?" == "0" ]; then
        rm **/bench_results*.json
        mkdir -p results
        if [[ -f /tmp/bench_results/result.json ]]; then
            mv /tmp/bench_results/result.json results/bench_final_results.$(date +%s).json
        fi
    else
        testErrorOut "Generate report failed"
        exit 1
    fi
}

function dcBuild() {
    docker-compose build
}

function dcRun() {
    dcUp "No Volumes" "docker-compose.yml"
    dcCollect "bench_results.json"
    dcDown "docker-compose.yml"

    dcUp "Local Disk Volume" "docker-compose.disk.yml"
    dcCollect "bench_results.disk.json"
    dcDown "docker-compose.disk.yml"

    dcUp "Remote TCP Volume" "docker-compose.guss.yml"
    dcCollect "bench_results.guss.json"
    dcDown "docker-compose.guss.yml"

    dcUp "Remote DMA Volume" "docker-compose.mellanox.yml"
    dcCollect "bench_results.mellanox.json"
    dcDown "docker-compose.mellanox.yml"
}

function dcVerify() {
    docker ps >> /dev/null
    if [ "$?" != "0" ]; then
        testErrorOut "Could not connect to docker!? (docker ps failed)"
        exit 1
    fi
}

function dcGo() {
    testOut "Starting benchmark..."

    dcVerify

    dcBuild
    dcRun

    testOut "benchmark Finished."
}

cd "$(dirname $0)"

mkdir -p /tmp/bench_results
# Iterate over tests
ran="0"
for i in $(ls -d */ | grep -v results | grep ".*$1.*"); do
    ran="1"
    testOut "Testing $i"
    cd $i
    dcGo
    cd ..
done

if [ "$ran" == "1" ]; then
    dcCollectAll
fi
