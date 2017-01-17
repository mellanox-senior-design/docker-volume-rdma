package driver

import (
	"github.com/golang/glog"
	"github.com/docker/go-plugins-helpers/volume"
)

// Holds all the information pertaining to a RDMA Volume Driver.
type RDMAVolumeDriver struct {

}

func NewRDMAVolumeDriver() RDMAVolumeDriver {
	return RDMAVolumeDriver{}
}

//	Create a new volume with name and options.
//	POST /VolumeDriver.Create
// 		in: { "Name": "volume_name", "Opts": {} }
//		out: { "Err": "" }
// 			Respond with a string error if an error occurred.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Create(request volume.Request) volume.Response {
	glog.Info("Creating volume: " + request.Name)

	return volume.Response{}
}

//	Get the list of volumes registered with the plugin.
//	POST /VolumeDriver.List
// 		in: {}
//		out: { "Volumes": [ { "Name": "volume_name", "Mountpoint": "/path/to/directory/on/host" } ], "Err": "" }
// 				Respond with a string error if an error occurred. Mountpoint is optional.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) List(request volume.Request) volume.Response {
	glog.Info("Listing volumes")

	return volume.Response{}
}

//	Get info relating to a paricular volume.
//	POST /VolumeDriver.Get
// 		in: { "Name": "volume_name" }
//		out: { "Volume": { "Name": "volume_name", "Mountpoint": "/path/to/directory/on/host", "Status": {} }, "Err": "" }
// 				Respond with a string error if an error occurred. Mountpoint and Status are optional.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Get(request volume.Request) volume.Response {
	glog.Info("Getting volume: " + request.Name)

	return volume.Response{}
}

//	Remove/Delete a paricular volume.
//	POST /VolumeDriver.Remove
// 		in: { "Name": "volume_name" }
//		out: { "Err": "" }
// 				Respond with a string error if an error occurred.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Remove(request volume.Request) volume.Response {
	glog.Info("Removing volume: " + request.Name)

	return volume.Response{}
}

//	Remind docker of the path a particular volume is attached to.
//	POST /VolumeDriver.Path
// 		in: { "Name": "volume_name" }
//		out: { "Mountpoint": "/path/to/directory/on/host", "Err": "" }
// 				Respond with the path on the host filesystem where the volume
//				has been made available, and/or a string error if an error
//				occurred. Mountpoint is optional, however the plugin may be queried
//				again later if one is not provided.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Path(request volume.Request) volume.Response {
	glog.Info("Getting path of volume: " + request.Name)

	return volume.Response{}
}

//	Mount a paricular volume to the host.
//	POST /VolumeDriver.Mount
// 		in: { "Name": "volume_name", "ID": "b87d7442095999a92b65b3d9691e697b61713829cc0ffd1bb72e4ccd51aa4d6c" }
//				Docker requires the plugin to provide a volume, given a user specified
//				volume name. This is called once per container start. If the same
//				volume_name is requested more than once, the plugin may need to keep
//				track of each new mount request and provision at the first mount request
//				and deprovision at the last corresponding unmount request.
//				ID is a unique ID for the caller that is requesting the mount.
//		out: { "Mountpoint": "/path/to/directory/on/host", "Err": "" }
// 				Respond with the path on the host filesystem where the volume
//				has been made available, and/or a string error if an error occurred.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Mount(request volume.MountRequest) volume.Response {
	glog.Info("Mounting volume: " + request.Name)

	return volume.Response{}
}

//	Remove a paricular volume from host (if no other mount requests active).
//	POST /VolumeDriver.Unmount
// 		in: { "Name": "volume_name", "ID": "b87d7442095999a92b65b3d9691e697b61713829cc0ffd1bb72e4ccd51aa4d6c" }
//				Indication that Docker no longer is using the named volume. This
//				is called once per container stop. Plugin may deduce that it is
//				safe to deprovision it at this point.
//				ID is a unique ID for the caller that is requesting the mount.
//		out: { "Err": "" }
// 				Respond with a string error if an error occurred.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Unmount(request volume.UnmountRequest) volume.Response {
	glog.Info("Unmounting volume: " + request.Name)

	return volume.Response{}
}

//	List the Docker Plugin capabilities that our plugin supports.
//	POST /VolumeDriver.Capabilities
// 		in: {}
//				Get the list of capabilities the driver supports. The driver is
//				not required to implement this endpoint, however in such cases
//				the default values will be taken.
//		out: { "Capabilities": { "Scope": "global" } }
//				Supported scopes are global and local. Any other value in Scope
//				will be ignored and assumed to be local. Scope allows cluster
//				managers to handle the volume differently, for instance with a
//				scope of global, the cluster manager knows it only needs to
//				create the volume once instead of on every engine. More
//				capabilities may be added in the future.
//		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Capabilities(request volume.Request) volume.Response {
	glog.Info("Listing capabilities")

	var response volume.Response
	response.Capabilities = volume.Capability{Scope: "local"}
	return response
}
