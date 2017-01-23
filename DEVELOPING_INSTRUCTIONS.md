# How to develop project

## Getting Started

### Install Golang

MacOS with [Brew][brew]

```bash
brew install go
```
Debian/Ubuntu with apt

```bash
apt install golang
```
> This will require root/sudo access


### Update `$PATH`

Golang requires that your go development happens in the same folder: your `$GOPATH`. We will need to add an export statement to your profile file. This varies for each shell, for brevity we have included the one shell `bash`. You will need to modify the instructions if you use other shells.

For `bash`, update your `~/.bashrc`
```bash
echo 'export GOPATH="/path/to/your/go/projects"' >> ~/.bashrc
echo 'export PATH="$PATH:$GOPATH/bin"' >> ~/.bashrc
source ~/.bashrc
echo $GOPATH
```

> Please note that every go project must be a subdirectory to your `$GOPATH`. This project must be a subdirectory to your `$GOPATH`.

#### Hello, world!

If you wish to test the installation and configuration of Golang, we can create a sample hello world application.

*Create a folder named `hello` in `$GOPATH`*
```bash
mkdir -p $GOPATH/src/hello && cd $GOPATH/src/hello
```
*Create a file named `main.go` and edit it*
```bash
vi main.go
```

Set the contents to:
```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
}
```

*Run `main.go`*
```bash
go build
./hello

# or one line
go run main.go
```

### Download REPO

*Checking out*
```bash
mkdir -p $GOPATH/src/github.com/mellanox-senior-design
cd $GOPATH/src/github.com/mellanox-senior-design
git clone git@github.com:mellanox-senior-design/docker-volume-rdma.git
```

*Downloading required libraries*
```bash
go get ./...
```

## Running the Volume Driver

If you are running the code from a machine that is not running a docker engine locally, you will need to make the `/etc/docker` folder.
```bash
sudo mkdir -p /etc/docker
sudo chmod -R 777 /etc/docker
```

*Start the server*
```bash
./run.sh

# or manually
cd $GOPATH/src/github.com/mellanox-senior-design # This is here incase that you are in a subdirectory or have a symlink to the $GOPATH folder
go run main.go -logtostderr=true
```

*Access the server*
```bash
curl -X "POST" "http://localhost:8080/VolumeDriver.Create" \
     -H "Content-Type: application/json" \
     -d $'{"Name": "volume_name", "Opts": {}}'
```
> NOTE: All endpoints use the POST method!
Accessing the server via a web browser will never return a nice result.

[brew]: http://brew.sh
