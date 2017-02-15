#!/bin/bash

docker run \
    --detach \
    --name jenkins \
    --publish 80:8080/tcp \
    --volume /var/run/docker.sock:/var/run/docker.sock \
    --volume /var/data/jenkins:/var/jenkins \
    onesysadmin/jenkins-docker-executors:latest
