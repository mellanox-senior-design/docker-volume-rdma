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
We have two options here. Either we can checkout the code directly into the `$GOPATH/src` folder or we can create a symlink to the code.

*Checking out*
```bash
cd $GOPATH/src #or the path you actually want to put this
git clone git@github.com:Jacobingalls/EE464K-RDMA-docker-volume-server.git
```

*Creating symlink*
```bash
# If you did not put the project in your $GOPATH, create a symbolic link to it.
if [ ! -d "$GOPATH/src/E464K-RDMA-docker-volume-server" ]; then
    ln -s E464K-RDMA-docker-volume-server $GOPATH/src/E464K-RDMA-docker-volume-server
fi
```

[brew]: http://brew.sh
