package drivers

import (
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
	"github.com/mellanox-senior-design/docker-volume-rdma/db"
)

// RDMAVolumeDriver holds all the information pertaining to a RDMA Volume Driver.
type RDMAVolumeDriver struct {
	StorageController RDMAStorageController
	VolumeDatabase    db.VolumeDatabase
}

// RDMAStorageController interface allowing Storage Controllers to create, mounte, remove, ect. volumes on a host.
type RDMAStorageController interface {
	Connect() error
	Disconnect() error

	// Mount a volume and mark that ID is using it, returning Mountpoint and error (nil if no error).
	Mount(volumeName string, id string) (string, error)

	// Unmount a volume from the host (NOT DELETE) if there are no longer any other ID's asking for the volume, returning error (nil if no error).
	Unmount(volumeName string, id string) error
}

// NewRDMAVolumeDriver constructs a new RDMAVolumeDriver.
func NewRDMAVolumeDriver(sc RDMAStorageController, vd db.VolumeDatabase) RDMAVolumeDriver {
	return RDMAVolumeDriver{sc, vd}
}

func (r RDMAVolumeDriver) validateOrCrash() {
	if r.StorageController == nil {
		glog.Fatal("StorageController is nil! Please configure the StorageController in the RDMAVolumeDriver.")
	}
}

// Connect to both the volume driver and the storage controller
func (r RDMAVolumeDriver) Connect() error {
	var err error
	err = r.VolumeDatabase.Connect()
	if err != nil {
		return err
	}

	err = r.StorageController.Connect()
	if err != nil {
		return err
	}

	return nil
}

// Disconnect from both the volume driver and the storage controller
func (r RDMAVolumeDriver) Disconnect() error {
	var err error
	err = r.VolumeDatabase.Disconnect()
	if err != nil {
		return err
	}

	err = r.StorageController.Disconnect()
	if err != nil {
		return err
	}

	return nil
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

	// Pass the create request to the volume database.
	err := r.VolumeDatabase.Create(request.Name, request.Options)

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

	// Pass the list request to the volume database.
	vols, err := r.VolumeDatabase.List()

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

	// Pass the get request to the volume database.
	vol, err := r.VolumeDatabase.Get(request.Name)

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

	// TODO unmount all mounts before remove

	// Pass the remove request to the volume database.
	err := r.VolumeDatabase.Remove(request.Name)

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
	mountpoint, err := r.VolumeDatabase.Path(request.Name)

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

	// Pass the mount request to the volume database.
	mountpoint, err := r.VolumeDatabase.Mount(request.Name, request.ID)

	// Pass the mount request to the storage controller.
	// TODO mountpoint, err := r.StorageController.Mount(request.Name, request.ID)

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

	// Pass the unmount request to the volume database.
	err := r.VolumeDatabase.Unmount(request.Name, request.ID)

	// Pass the unmount request to the storage controller.
	// TODO err := r.StorageController.Unmount(request.Name, request.ID)

	// If there was an error, log.
	if err != nil {
		errString = err.Error()
		glog.Error("Error: " + errString + " Encountered while unmounting volume: " + request.Name + " and ID: " + request.ID)
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
