#! /bin/bash

cd ubuntu

docker build -t bonnie .
docker run -it --rm bonnie

cd ..

cd debian

docker build -t bonnie .
docker run -it --rm bonnie

cd ..
