package db

import "github.com/docker/go-plugins-helpers/volume"

type VolumeDatabase interface {
	Connect() error
	Disconnect() error

	Create(volumeName string, options map[string]string) error
	List() ([]*volume.Volume, error)
	Get(volumeName string) (*volume.Volume, error)
	Path(volumeName string) (string, error)
	Remove(volumeName string) error

	Mount(volumeName string, id string) (string, error)
	Unmount(volumeName string, id string) error
}
