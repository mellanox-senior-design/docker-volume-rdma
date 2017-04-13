#! /bin/bash



go test -cover github.com/mellanox-senior-design/docker-volume-rdma/...

for i in $(ls *.go); do
lines=$(cat $i | grep -E "[a-zA-Z]" | grep -vE "^\s*\/\/.*" | grep -v func | wc -l)
echo "$i $lines"
done

for i in $(ls **/*.go); do
lines=$(cat $i | grep -E "[a-zA-Z]" | grep -vE "^\s*\/\/.*" | grep -v func | wc -l)
echo "$i $lines"
done
