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
    docker-compose down -v
}

function dcUp() {
    testOut "Bringing up Test Fixture [$1]..."
    if [ -f $2 ]; then
        docker-compose up --abort-on-container-exit --force-recreate $file
    else
        testWarnOut "File, $2, does not exist. Skipping..."
    fi
}

function dcCollect() {
    testOut "Collecting..."
    mv /tmp/bench_results/result.json "$1"
    testOut "$(cat $1)"
}

function dcCollectAll() {
    testOut "Collecting All..."
    ./generate_report.py > /tmp/bench_results/result.json
    if [ "$?" == "0" ]; then
        rm **/bench_results*.json
        mkdir -p results
        mv /tmp/bench_results/result.json results/bench_final_results.$(date +%s).json
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
    dcDown

    dcUp "Local Disk Volume" "docker-compose.disk.yml"
    dcCollect "bench_results.disk.json"
    dcDown

    dcUp "Remote DMA Volume" "docker-compose.rdma.yml"
    dcCollect "bench_results.rdma.json"
    dcDown
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
    dcDown

    testOut "benchmark Finished."
}

# Itterate over tests
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
