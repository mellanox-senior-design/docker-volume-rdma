package drivers

// GlusterStorageController connects to the local gluster client and facilitates volume mounts
type GlusterStorageController struct {
}

// NewGlusterStorageController creates a new GlusterStorageController
func NewGlusterStorageController() GlusterStorageController {
	return GlusterStorageController{}
}

// Connect to glusterfs
func (g GlusterStorageController) Connect() error {
	return nil
}

// Disconnect from glusterfs
func (g GlusterStorageController) Disconnect() error {
	return nil
}

// Mount a volume by name
func (g GlusterStorageController) Mount(volumeName string) (string, error) {
	return "", nil
}

// Unmount a volume by name
func (g GlusterStorageController) Unmount(volumeName string) error {
	return nil
}

// Delete a volume, permanatly remove data
func (g GlusterStorageController) Delete(volumeName string) error {
	return nil
}
