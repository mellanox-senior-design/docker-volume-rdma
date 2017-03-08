#! /bin/bash

while ! nc -z mysql 3306; do
  echo "Waiting for MySQL to launch (3306)."
  sleep 1
done
echo "MySQL is up!"

sleep 5

echo "CREATE SCHEMA rdma;" | mysql -u root -ppassword -h mysql
echo $?

./docker-volume-rdma -logtostderr=true "$@"
