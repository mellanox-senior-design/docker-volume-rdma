package main

import (
	"fmt"
	"log"
	"github.com/docker/go-plugins-helpers/volume"
)

type RDMAVolumeDriver struct {

}

func (r RDMAVolumeDriver) Create(request volume.Request) volume.Response {
	log.Println("Creating volume: " + request.Name)

	return volume.Response{}
}

func (r RDMAVolumeDriver) List(request volume.Request) volume.Response {
	log.Println("Listing volumes")

	return volume.Response{}
}

func (r RDMAVolumeDriver) Get(request volume.Request) volume.Response {
	log.Println("Getting volume: " + request.Name)

	return volume.Response{}
}

func (r RDMAVolumeDriver) Remove(request volume.Request) volume.Response {
	log.Println("Removing volume: " + request.Name)

	return volume.Response{}
}

func (r RDMAVolumeDriver) Path(request volume.Request) volume.Response {
	log.Println("Getting path of volume: " + request.Name)

	return volume.Response{}
}

func (r RDMAVolumeDriver) Mount(request volume.MountRequest) volume.Response {
	log.Println("Mounting volume: " + request.Name)

	return volume.Response{}
}

func (r RDMAVolumeDriver) Unmount(request volume.UnmountRequest) volume.Response {
	log.Println("Unmounting volume: " + request.Name)

	return volume.Response{}
}

func (r RDMAVolumeDriver) Capabilities(request volume.Request) volume.Response {
	log.Println("Listing capabilities")

	var response volume.Response
    response.Capabilities = volume.Capability{Scope: "local"}
    return response
}

func main() {
	fmt.Println("Starting...")

	driver := RDMAVolumeDriver{}
	h := volume.NewHandler(driver)
	err := h.ServeTCP("test_volume", ":8080", nil)

	fmt.Println(err)
}
