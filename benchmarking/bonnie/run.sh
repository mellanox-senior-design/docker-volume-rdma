#! /bin/bash

bonnie++ -d /test -r 512 -u root | grep $(hostname) | head -n 1 >> result.txt
python bench.py
