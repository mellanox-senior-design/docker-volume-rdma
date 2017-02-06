package db

import (
	"errors"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
)

// InMemoryVolumeDatabase defines a volume database that is completely in memory (and ephimeral)
type InMemoryVolumeDatabase struct {
	volumes map[string]*volume.Volume
	mounts  map[string]map[string]int
}

// NewInMemoryVolumeDatabase creates a new InMemoryVolumeDatabase, inilizing all of its properties.
func NewInMemoryVolumeDatabase() InMemoryVolumeDatabase {
	warning := `
********************************************************************************
*                                                                              *
*                                WARNING!!!                                    *
*              DO NOT USE the "In-Memory Driver" IN PRODUCTION!                *
*                                                                              *
*        You are currently using the In-Memory Volume Database Driver,         *
*        it is completely ephemeral and should only be used for testing        *
*        purposes. When the plugin process stops, all data will be lost.       *
*                                                                              *
*        However, any data saved in the volumes will still be available        *
*        by manually accessing the data in the storage backend filesystem.     *
*                                                                              *
********************************************************************************
	`

	glog.Warning(warning)

	volumes := map[string]*volume.Volume{}
	mounts := map[string]map[string]int{}
	return InMemoryVolumeDatabase{volumes: volumes, mounts: mounts}
}

// Connect is a NOP, though required by VolumeDatabase interface
func (i InMemoryVolumeDatabase) Connect() error {
	glog.Info("Connect function called, no action taken.")
	return nil
}

// Disconnect is a NOP, though required by VolumeDatabase interface
func (i InMemoryVolumeDatabase) Disconnect() error {
	glog.Info("Disconnect function called, no action taken.")
	return nil
}

// Create a new volume in the database, returning an error if one occured.
func (i InMemoryVolumeDatabase) Create(volumeName string, options map[string]string) error {
	var exists bool
	_, exists = i.volumes[volumeName]
	if exists {
		return errors.New("Volume already exists.")
	}

	_, exists = i.mounts[volumeName]
	if exists {
		return errors.New("Volume already exists.")
	}

	i.volumes[volumeName] = &volume.Volume{
		Name:       volumeName,
		Mountpoint: "",
		Status:     nil}

	return nil
}

// List all of the volumes in the database, returning an error if one occured.
func (i InMemoryVolumeDatabase) List() ([]*volume.Volume, error) {
	volumeList := make([]*volume.Volume, 0, len(i.volumes))

	for _, value := range i.volumes {
		volumeList = append(volumeList, value)
	}

	return volumeList, nil
}

// Get the volume definiton from the database by name, returning an error if one occured.
func (i InMemoryVolumeDatabase) Get(volumeName string) (*volume.Volume, error) {
	vol, exists := i.volumes[volumeName]
	if !exists {
		return nil, errors.New("Volume does not exist.")
	}

	return vol, nil
}

// Path of the specified volume, returning an error if one occured.
func (i InMemoryVolumeDatabase) Path(volumeName string) (string, error) {
	vol, err := i.Get(volumeName)
	if err != nil {
		return "", err
	}

	return vol.Mountpoint, nil
}

// Remove (delete) the specified volume, returning an error if one occured.
func (i InMemoryVolumeDatabase) Remove(volumeName string) error {
	var exists bool
	_, exists = i.volumes[volumeName]
	if !exists {
		return errors.New("Volume does not exit.")
	}
	delete(i.volumes, volumeName)

	_, exists = i.mounts[volumeName]
	if exists {
		sum := 0
		for _, v := range i.mounts[volumeName] {
			sum += v
		}

		if sum != 0 {
			return errors.New("Cannot remove volume as it is still being requested.")
		}
	}

	delete(i.mounts, volumeName)
	return nil
}

// Mount the specified volume to the host and increment the id to prevent premature removal, returning an error if one occured.
func (i InMemoryVolumeDatabase) Mount(volumeName string, id string) (string, error) {
	vol, err := i.Get(volumeName)
	if err != nil {
		return "", err
	}

	var exists bool
	_, exists = i.mounts[volumeName]
	if exists {
		_, exists = i.mounts[volumeName][id]
		if exists {
			i.mounts[volumeName][id]++
		} else {
			i.mounts[volumeName][id] = 1
		}
	} else {
		i.mounts[volumeName] = map[string]int{}
		i.mounts[volumeName][id] = 1
	}

	mountpoint := "/fake/location/" + volumeName
	vol.Mountpoint = mountpoint
	glog.Info(volumeName, " and id ", id, " is now has ", i.mounts[volumeName][id], " connections to "+mountpoint)
	return mountpoint, nil
}

// Unmount the specified volume if the ids are no longer referencing it, returning an error if one occured.
func (i InMemoryVolumeDatabase) Unmount(volumeName string, id string) error {
	_, err := i.Get(volumeName)
	if err != nil {
		return err
	}

	var exists bool
	_, exists = i.mounts[volumeName]
	if exists {
		_, exists = i.mounts[volumeName][id]
		if exists {
			if i.mounts[volumeName][id] > 0 {
				i.mounts[volumeName][id]--
			} else {
				return errors.New("Volume exists, but was not mounted.")
			}
		} else {
			i.mounts[volumeName][id] = 0
			return errors.New("Volume exists, but was not mounted.")
		}
	} else {
		i.mounts[volumeName] = map[string]int{}
		i.mounts[volumeName][id] = 0
		return errors.New("Volume exists, but was not mounted.")
	}

	glog.Info(volumeName, " and id ", id, " is now has ", i.mounts[volumeName][id], " connections")
	return nil
}
