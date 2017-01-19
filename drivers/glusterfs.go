package drivers

type GlusterStorageController struct {
}

func NewGlusterStorageController() GlusterStorageController {
	return GlusterStorageController{}
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
