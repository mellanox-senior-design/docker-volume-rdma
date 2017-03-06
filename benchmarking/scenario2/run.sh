# #! /bin/bash -x


docker-compose up -d

docker build -t hey .
docker run -it --rm hey

docker-compose down

# python getblogpost.py 1 100
