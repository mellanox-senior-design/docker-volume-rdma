package drivers

import (
	"errors"
	"os"
	"path"

	"github.com/golang/glog"
)

// OnDiskStorageController is a way of testing the StorageController backend to ensure that all of the logic is correct.
type OnDiskStorageController struct {
	FSPath string
}

// NewOnDiskStorageController creates a new OnDiskStorageController
func NewOnDiskStorageController(path string) OnDiskStorageController {

	warning := `
********************************************************************************
*                                                                              *
*                                WARNING!!!                                    *
*             DO NOT USE the "On-Disk Controller" IN PRODUCTION!               *
*                                                                              *
*        You are currently using the On-Disk Storage Controller, it            *
*        operates on the host's disk with out the benifits of rdma and         *
*        should only be used for testing purposes. All data saved in           *
*        the volumes will still be available by manually accessing the         *
*        data in the host's filesystem.                                        *
*                                                                              *
********************************************************************************
	`

	glog.Warning(warning)

	if path == "" {
		path = "/etc/docker/mounts/"
	}

	glog.Info("Mount path: ", path)

	_, err := os.Open(path)
	if err != nil {
		os.MkdirAll(path, 0755)
		_, err = os.Open(path)
		if err != nil {
			glog.Fatal("Unable to create folder ", path, ". ", err)
		}
	}

	return OnDiskStorageController{path}
}

// Connect is a NOOP
func (d OnDiskStorageController) Connect() error {
	glog.Info("Connect function called, no action taken.")
	return nil
}

// Disconnect is a NOOP
func (d OnDiskStorageController) Disconnect() error {
	glog.Info("Disconnect function called, no action taken.")
	return nil
}

// Mount a particular volume
func (d OnDiskStorageController) Mount(volumeName string) (string, error) {
	pathMounted := path.Join(d.FSPath, volumeName)
	pathUnmounted := path.Join(path.Dir(pathMounted), path.Base(pathMounted)+".unmounted")

	// If there is an unmounted volume, return it.
	_, err := os.Open(pathUnmounted)
	if err == nil {
		glog.Info("Renaming: ", pathMounted, " to ", pathUnmounted)
		return pathMounted, os.Rename(pathUnmounted, pathMounted)
	}

	_, err = os.Open(pathMounted)
	if err != nil {
		glog.Info("Creating: ", pathMounted)
		os.MkdirAll(pathMounted, 0755)
		_, err = os.Open(pathMounted)
		if err != nil {
			return "", err
		}
	}

	return pathMounted, nil
}

// Unmount a particular volume
func (d OnDiskStorageController) Unmount(volumeName string) error {
	pathMounted := path.Join(d.FSPath, volumeName)
	pathUnmounted := path.Join(path.Dir(pathMounted), path.Base(pathMounted)+".unmounted")

	glog.Info(pathMounted)

	// If there is an unmounted volume, return it.
	_, err := os.Open(pathMounted)
	if err == nil {
		glog.Info("Renaming: ", pathMounted, " to ", pathUnmounted)
		return os.Rename(pathMounted, pathUnmounted)
	}

	_, err = os.Open(pathUnmounted)
	if err != nil {
		return err
	}

	return errors.New("already unmounted")
}

// Delete a particular volume
func (d OnDiskStorageController) Delete(volumeName string) error {
	pathMounted := path.Join(d.FSPath, volumeName)
	pathUnmounted := path.Join(path.Dir(pathMounted), path.Base(pathMounted)+".unmounted")

	// If there is an unmounted volume, return it.
	_, err := os.Open(pathUnmounted)
	if err == nil {
		return os.Remove(pathUnmounted)
	}

	_, err = os.Open(pathMounted)
	if err == nil {
		return os.Remove(pathMounted)
	}

	return nil
}
