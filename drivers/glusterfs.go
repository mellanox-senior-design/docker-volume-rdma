package drivers

import (
	"errors"

	"github.com/docker/go-plugins-helpers/volume"
)

type GlusterStorageController struct {
}

func NewGlusterStorageController() GlusterStorageController {
	return GlusterStorageController{}
}

func (g GlusterStorageController) Create(volumeName string, opts map[string]string) error {
	return nil
}

func (g GlusterStorageController) List() ([]*volume.Volume, error) {
	return nil, nil
}

func (g GlusterStorageController) Get(volumeName string) (*volume.Volume, error) {
	return nil, nil
}

func (g GlusterStorageController) Remove(volumeName string) error {
	return nil
}

func (g GlusterStorageController) Path(volumeName string) (string, error) {
	return "", nil
}

func (g GlusterStorageController) Mount(volumeName string, id string) (string, error) {
	return "/some/path/here", nil
}

func (g GlusterStorageController) Unmount(volumeName string, id string) error {
	return nil
}
