package db

import "github.com/docker/go-plugins-helpers/volume"

// VolumeDatabase interface describes a connection to a database that is capable of keeping track of volumes and mounts
type VolumeDatabase interface {

	// Connect to the database, and create anything needed to process requests.
	// note: this function is run on the main thread, not a goroutine.
	Connect() error

	// Disconnect from the database, close any connections etc.
	Disconnect() error

	// Create volume by name and options.
	Create(volumeName string, options map[string]string) error

	// List all of the volumes that we know about.
	List() ([]*volume.Volume, error)

	// Get info about a particular volume.
	Get(volumeName string) (*volume.Volume, error)

	// Get the path of a particular volume.
	Path(volumeName string) (string, error)

	// Remove a particular volume.
	Remove(volumeName string) error

	// Mount a particular volume.
	Mount(volumeName string, id string, mountpoint string) error

	// Unmount a particular volume.
	Unmount(volumeName string, id string) error
}
