#! /bin/bash

#  Wait for WordPress blog
#  http://stackoverflow.com/a/27601038/3259030
while ! nc -z mysql 3306; do
  echo "Waiting for MySQL to launch (3306)."
  sleep 1
done

echo "MySql is up!"

sleep 5
flyway-4.1.1/flyway -user=root -password=wordpress -url=jdbc:mysql://mysql -schemas=wordpress migrate

while ! nc -z wordpress 80; do
  echo "Waiting for WordPress to launch (80)."
  sleep 1
done

echo "WordPress blog is up!"

python bench.py
