#! /bin/bash

sleep 1
flyway-4.1.1/flyway -user=root -password=password -url=jdbc:mysql://mysql -schemas=bench migrate
python bench.py
