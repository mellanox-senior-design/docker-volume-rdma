package drivers

type GlusterStorageController struct {
}

func NewGlusterStorageController() GlusterStorageController {
	return GlusterStorageController{}
}

func (g GlusterStorageController) Connect() error {
	return nil
}

func (g GlusterStorageController) Disconnect() error {
	return nil
}

func (g GlusterStorageController) Mount(volumeName string) (string, error) {
	return "", nil
}

func (g GlusterStorageController) Unmount(volumeName string) error {
	return nil
}

func (g GlusterStorageController) Delete(volumeName string) error {
	return nil
}
