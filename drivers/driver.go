package drivers

import (
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
)

// RDMAVolumeDriver holds all the information pertaining to a RDMA Volume Driver.
type RDMAVolumeDriver struct {
	StorageController RDMAStorageController
}

// RDMAStorageController interface allowing Storage Controllers to create, mounte, remove, ect. volumes on a host.
type RDMAStorageController interface {

	// Create a volume with name and opts, returning error (nil if no error).
	Create(volumeName string, opts map[string]string) error

	// Return the list of all volumes, returning the list of volumes and error (nil if no error).
	List() ([]*volume.Volume, error)

	// Get information about a volume by name, returning the volume and error (nil if no error).
	Get(volumeName string) (*volume.Volume, error)

	// Remove/Delete volume, returning an error (nil if no error).
	Remove(volumeName string) error

	// Get the mounted path of a volume, returning the Mountpoint and error (nil if no error).
	Path(volumeName string) (string, error)

	// Mount a volume and mark that ID is using it, returning Mountpoint and error (nil if no error).
	Mount(volumeName string, id string) (string, error)

	// Unmount a volume from the host (NOT DELETE) if there are no longer any other ID's asking for the volume, returning error (nil if no error).
	Unmount(volumeName string, id string) error
}

// NewRDMAVolumeDriver constructs a new RDMAVolumeDriver.
func NewRDMAVolumeDriver(sc RDMAStorageController) RDMAVolumeDriver {
	return RDMAVolumeDriver{sc}
}

func (r RDMAVolumeDriver) validateOrCrash() {
	if r.StorageController == nil {
		glog.Fatal("StorageController is nil! Please configure the StorageController in the RDMAVolumeDriver.")
	}
}

// Create a new volume with name and options.
// POST /VolumeDriver.Create
// 		in: { "Name": "volume_name", "Opts": {} }
//		out: { "Err": "" }
// 			Respond with a string error if an error occurred.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Create(request volume.Request) volume.Response {
	glog.Info("Creating volume: " + request.Name)

	// Ensure the r is properly configured
	r.validateOrCrash()

	// Pass the create request to the storage controller.
	err := r.StorageController.Create(request.Name, request.Options)

	// If there was an error, log.
	var errString string
	if err != nil {
		errString = err.Error()
		glog.Error("Error: " + errString + "! Encountered while creating a volume: " + request.Name)
	}

	// Construct and return a response using the docker library.
	var response volume.Response
	response.Err = errString
	return response
}

// List gets the list of volumes registered with the plugin.
// POST /VolumeDriver.List
// 		in: {}
//		out: { "Volumes": [ { "Name": "volume_name", "Mountpoint": "/path/to/directory/on/host" } ], "Err": "" }
// 				Respond with a string error if an error occurred. Mountpoint is optional.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) List(request volume.Request) volume.Response {
	glog.Info("Listing volumes")

	// Ensure the r is properly confiured
	r.validateOrCrash()

	// Pass the list request to the storage controller.
	vols, err := r.StorageController.List()

	// If there was an error, log.
	var errString string
	if err != nil {
		errString = err.Error()
		glog.Error("Error: " + errString + "! Encountered while listing the volumes")
	}

	// Construct and return a response using the docker library.
	var response volume.Response
	response.Volumes = vols
	response.Err = errString
	return response
}

// Get info relating to a paricular volume.
// POST /VolumeDriver.Get
// 		in: { "Name": "volume_name" }
//		out: { "Volume": { "Name": "volume_name", "Mountpoint": "/path/to/directory/on/host", "Status": {} }, "Err": "" }
// 				Respond with a string error if an error occurred. Mountpoint and Status are optional.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Get(request volume.Request) volume.Response {
	glog.Info("Getting volume: " + request.Name)

	// Ensure the r is properly confiured
	r.validateOrCrash()

	// Pass the get request to the storage controller.
	vol, err := r.StorageController.Get(request.Name)

	// If there was an error, log.
	var errString string
	if err != nil {
		errString = err.Error()
		glog.Error("Error: " + errString + "! Encountered while get information about the volume: " + request.Name)
	}

	// Construct and return a response using the docker library.
	var response volume.Response
	response.Volume = vol
	response.Err = errString
	return response
}

// Remove (Delete) a paricular volume.
// POST /VolumeDriver.Remove
// 		in: { "Name": "volume_name" }
//		out: { "Err": "" }
// 				Respond with a string error if an error occurred.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Remove(request volume.Request) volume.Response {
	glog.Info("Removing volume: " + request.Name)

	// Ensure the r is properly confiured
	r.validateOrCrash()

	// Pass the remove request to the storage controller.
	err := r.StorageController.Remove(request.Name)

	// If there was an error, log.
	var errString string
	if err != nil {
		errString = err.Error()
		glog.Error("Error: " + errString + "! Encountered while removing a volume: " + request.Name)
	}

	// Construct and return a response using the docker library.
	var response volume.Response
	response.Err = errString
	return response
}

// Path reminds docker of the path a particular volume is attached to.
// POST /VolumeDriver.Path
// 		in: { "Name": "volume_name" }
//		out: { "Mountpoint": "/path/to/directory/on/host", "Err": "" }
// 				Respond with the path on the host filesystem where the volume
//				has been made available, and/or a string error if an error
//				occurred. Mountpoint is optional, however the plugin may be queried
//				again later if one is not provided.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Path(request volume.Request) volume.Response {
	glog.Info("Getting path of volume: " + request.Name)

	// Ensure the r is properly confiured
	r.validateOrCrash()

	// Pass the path request to the storage controller.
	mountpoint, err := r.StorageController.Path(request.Name)

	// If there was an error, log.
	var errString string
	if err != nil {
		errString = err.Error()
		glog.Error("Error: " + errString + "! Encountered while get path of volume: " + request.Name)
	}

	// Construct and return a response using the docker library.
	var response volume.Response
	response.Mountpoint = mountpoint
	response.Err = errString
	return response
}

// Mount a paricular volume to the host.
// POST /VolumeDriver.Mount
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

	// Ensure the r is properly confiured
	r.validateOrCrash()

	// Pass the mount request to the storage controller.
	mountpoint, err := r.StorageController.Mount(request.Name, request.ID)

	// If there was an error, log.
	var errString string
	if err != nil {
		errString = err.Error()
		glog.Error("Error: " + errString + "! Encountered while mounting volume: " + request.Name + " and ID: " + request.ID)
	}

	// Construct and return a response using the docker library.
	var response volume.Response
	response.Mountpoint = mountpoint
	response.Err = errString
	return response
}

// Unmount a paricular volume from host (if no other mount requests active).
// POST /VolumeDriver.Unmount
// 		in: { "Name": "volume_name", "ID": "b87d7442095999a92b65b3d9691e697b61713829cc0ffd1bb72e4ccd51aa4d6c" }
//				Indication that Docker no longer is using the named volume. This
//				is called once per container stop. Plugin may deduce that it is
//				safe to deprovision it at this point.
//				ID is a unique ID for the caller that is requesting the mount.
//		out: { "Err": "" }
// 				Respond with a string error if an error occurred.
// 		referTo: https://docs.docker.com/engine/extend/plugins_volume/
func (r RDMAVolumeDriver) Unmount(request volume.UnmountRequest) volume.Response {
	var errString string
	glog.Info("Unmounting volume: " + request.Name + " and ID: " + request.ID)

	r.validateOrCrash()

	// Pass the unmount request to the storage controller.
	err := r.StorageController.Unmount(request.Name, request.ID)

	// If there was an error, log.
	if err != nil {
		errString = err.Error()
		glog.Error("Error: " + errString + "! Encountered while unmounting volume: " + request.Name + " and ID: " + request.ID)
	}

	// Construct and return a response using the docker library.
	var response volume.Response
	response.Err = errString
	return response
}

// Capabilities that our plugin supports.
// POST /VolumeDriver.Capabilities
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

	// Construct and return a response using the docker library.
	var response volume.Response
	response.Capabilities = volume.Capability{Scope: "local"}
	return response
}
