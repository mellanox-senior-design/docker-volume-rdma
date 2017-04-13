# RDMA Volume Plugin for Docker
[![Build Status](https://travis-ci.org/mellanox-senior-design/docker-volume-rdma.svg?branch=master)](https://travis-ci.org/mellanox-senior-design/docker-volume-rdma)

This volume plugin aims to add an RDMA-enabled storage backend to Docker Volumes.

> Capstone Project, 2017, [UT Austin][UT], Sponsored by [Mellanox][Mellanox]

[How to develop for docker-volume-plugin][develop]

[UT]: http://www.ece.utexas.edu
[develop]: DEVELOPING_INSTRUCTIONS.md
[Mellanox]: http://www.mellanox.com

## Connecting a container

```bash
docker run ... --volume-driver=docker-volume-rdma -v "volume_name:/path/to/put/volume:z" ...
```

As an example, lets launch a httpd server on a shared volume. On every connected
host we could run:

```bash
docker run -it --rm --volume-driver=docker-volume-rdma -v "volume_name:/usr/local/apache2/htdocs/:z" -p 80:80 httpd
```

Then we could launch another container, on a different machine, and edit the site.

```bash
docker run -it --rm --volume-driver=docker-volume-rdma -v "volume_name:/website:z" ubuntu vi /website/index.html
```

## Quick start

### Launch an instance of MySQL
Each host is going to run the plugin, so we require shared storage. On a machine
that has a port for mysql available, install mysql.

```bash
docker run -d \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD={{mysql_root_password}} \
    -v volume_db:/var/lib/mysql  \
    mysql
```

The database should be setup with a dedicated user, {{mysql_username}} and
password, {{mysql_password}}.

```mysql
CREATE SCHEMA {{mysql_schema}};
CREATE USER {{mysql_username}}@* IDENTIFIED BY {{mysql_password}};
GRANT ALL ON {{mysql_schema}}.* TO {{mysql_username}}@*;
```

### Connect storage solution
Connect the RDMA-enabled, shared storage to the container host. Mount it
to a common place on all machines, not strictly required, but it may be easier
to maintain.

For example: `/mnt/glusterfs/rdma/volumes`


### Install the plugin dependencies
Install Golang, and get the driver

```bash
# Ubuntu
apt-get install golang

# RHEL
yum install golang
```

Edit your ~/.bashrc or equivalent and configure the `GOPATH`, this will be where
all go projects store their dependencies.

```bash
export PATH=$PATH:$HOME/go    # Example, can be anything that you like better
```
> After editing, source the file again so that your `$PATH` is up to date.

### Install the plugin

```bash
go get github.com/mellanox-senior-design/docker-volume-rdma
```

### Run the driver in the background (TODO: make service)
```bash
cd $GOPATH/go/src/github.com/mellanox-senior-design/docker-volume-rdma
nohup ./run.sh \
    -db=mysql \
    -dbuser={{mysql_username}} \
    -dbpass={{mysql_password}} \
    -dbschema=rdma \
    -dbhost="tcp({{mysql}}:3306)" \
    -sc=on-disk \
    -scpath=/mnt/glusterfs/rmda/volumes &
```
