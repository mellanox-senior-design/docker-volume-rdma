#! /bin/bash

#  Wait for MySQL
#  http://stackoverflow.com/a/27601038/3259030
while ! nc -z mysql 3306; do
    echo "Waiting for MySQL to launch (3306)."
    sleep 1
done
echo "MySQL is up!"

sleep 5
flyway-4.1.1/flyway -user=root -password=password -url=jdbc:mysql://mysql -schemas=bench migrate
python bench.py
