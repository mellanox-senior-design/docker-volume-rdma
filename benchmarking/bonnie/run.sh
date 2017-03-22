#! /bin/bash

mkdir -p /test
bonnie++ -d /test -r 2048 -u root | grep $(hostname) | head -n 1 >> result.txt
python bench.py
