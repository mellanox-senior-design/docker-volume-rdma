// Creates a webserver/unix socket that allows a Docker server to create RDMA
// backed volumes.
package main

import (
	"flag"
	"strconv"

	"github.com/Jacobingalls/docker-volume-rdma/driver"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
)

// Port to launch service on.
var httpPort int

// Configure and start the docker volume plugin server.
func main() {

	// Configure application flags.
	flag.IntVar(&httpPort, "port", 8080, "tcp/ip port to serve volume driver on")

	// Parse flags as glog needs the flags to be solifified before starting.
	flag.Parse()

	// Convert port to string, and print startup message.
	port := strconv.Itoa(httpPort);
	glog.Info("Running! http://localhost:" + port)

	// Create and begin serving volume driver on tcp/ip port, httpPort.
	driver := driver.NewRDMAVolumeDriver()
	h := volume.NewHandler(driver)
	err := h.ServeTCP("test_volume", ":" + port, nil)

	// Log any error that may have occured.
	glog.Fatal(err)
}
