package drivers

import (
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/mellanox-senior-design/docker-volume-rdma/db"
)

/*
 * based on InMemDatabase & OnDiskStorageController
 * the driver will always connect correctly
 */
func TestConnect(t *testing.T) {
	t.Parallel()
	dataBase := db.NewInMemoryVolumeDatabase()
	storageController := NewOnDiskStorageController("")
	rdmaVolumeDriver := NewRDMAVolumeDriver(storageController, dataBase)

	err := rdmaVolumeDriver.Connect()
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidation(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)
	rdmaVolDriver.validateOrCrash()
}

func TestDisconnect(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)
	err := rdmaVolDriver.Disconnect()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)

	request := volume.Request{Name: "test"}
	response := rdmaVolDriver.Create(request)

	if len(response.Err) != 0 {
		t.Fatal(response.Err)
	}

	_, err := db.Get("test")

	if err != nil {
		t.Fatal(err)
	}

	response = rdmaVolDriver.Create(request)

	if len(response.Err) == 0 {
		t.Fatal("We should receive an error because a volume cannot be created twice")
	}

}

func TestList(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)
	req := volume.Request{Name: "testDos"}

	list := rdmaVolDriver.List(req)

	if len(list.Volumes) != 0 {
		t.Fatal("We should not have any volumes Listed as there were none created")
	}

	rdmaVolDriver.Create(req)

	list = rdmaVolDriver.List(req)

	if len(list.Err) != 0 {
		t.Fatal(list.Err)
	}

	if len(list.Volumes) != 1 {
		t.Fatal("Expected to see only one volume, but didn't. Saw ", len(list.Volumes), " number of volumes")
	}

	secondReq := volume.Request{Name: "testDosUno"}
	rdmaVolDriver.Create(secondReq)

	list = rdmaVolDriver.List(volume.Request{})

	if len(list.Err) != 0 {
		t.Fatal(list.Err)
	}

	if len(list.Volumes) != 2 {
		t.Fatal("Expected to see two volumes, but saw ", len(list.Volumes), " volumes")
	}

	for i := 0; i < len(list.Volumes); i++ {
		name := list.Volumes[i].Name
		if name != "testDos" {
			if name != "testDosUno" {
				t.Fatal("List grabbed the wrong volume. Expected testDos or testDosUno, but got "+list.Volumes[i].Name, i)
			}
		}
	}
}

func TestGet(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)

	response := rdmaVolDriver.Get(volume.Request{})

	if len(response.Err) == 0 {
		t.Fatal("Volume does not exist so there should be an error, but there is not")
	}

	if len(response.Volumes) != 0 {
		t.Fatal("Get should not return any volumes as there haven't been any created")
	}

	req := volume.Request{Name: "testGet"}

	rdmaVolDriver.Create(req)
	response = rdmaVolDriver.Get(req)

	if response.Volume.Name != "testGet" {
		t.Fatal("Wrong volume returned when calling Get")
	}

	if len(response.Err) != 0 {
		t.Fatal(response.Err)
	}

	secondReq := volume.Request{Name: "notCreated"}
	secondRes := rdmaVolDriver.Get(secondReq)
	if len(secondRes.Err) == 0 {
		t.Fatal("There should have been an error when Getting a volume that has not been created")
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)

	request := volume.Request{}

	response := rdmaVolDriver.Remove(request)

	if len(response.Err) == 0 {
		t.Fatal("Can't remove something that does not exist, we expect to receive an error")
	}

	request = volume.Request{Name: "removeMe"}

	rdmaVolDriver.Create(request)

	rdmaVolDriver.Create(volume.Request{Name: "keepMe"})

	response = rdmaVolDriver.Remove(request)

	if len(response.Err) != 0 {
		t.Fatal(response.Err)
	}

	_, err := db.Get("removeMe")

	if err == nil {
		t.Fatal("Volume should not exist, so we should not get one back")
	}

	_, err = db.Get("keepMe")

	if err != nil {
		t.Fatal(err)
	}
}

func TestPath(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)

	request := volume.Request{Name: "vol1"}

	response := rdmaVolDriver.Path(request)

	if len(response.Err) == 0 {
		t.Fatal("We should not hava a path sine volume was not created")
	}

	rdmaVolDriver.Create(request)
	response = rdmaVolDriver.Path(request)

	if len(response.Mountpoint) != 0 {
		t.Fatal("We should not have a path since we have not mounted the volume")
	}

	response = rdmaVolDriver.Mount(volume.MountRequest{Name: "vol1"})

	if len(response.Err) != 0 {
		t.Fatal("Encountered issue while mounting volume ", response.Err)
	}

	response = rdmaVolDriver.Path(request)

	if len(response.Err) != 0 {
		t.Fatal("Encountered issue while requesting path of volume ", response.Err)
	}

	if response.Mountpoint != "/etc/docker/mounts/vol1" {
		t.Fatal("Did not receive the expected path, instead got ", response.Mountpoint)
	}

}

func TestMount(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)

	request := volume.Request{Name: "movies"}
	// sRequest := volume.Request{Name: "songs"}
	mountRequest := volume.MountRequest{Name: "", ID: "42"}

	response := rdmaVolDriver.Mount(mountRequest)
	if len(response.Err) == 0 {
		t.Fatal("We should not be able to mount a volume whose name we don't know")
	}

	response = rdmaVolDriver.Create(request)

	if len(response.Err) != 0 {
		t.Fatal(response.Err)
	}

	mountRequest.Name = "movies"
	response = rdmaVolDriver.Mount(mountRequest)

	if len(response.Err) != 0 {
		t.Fatal(response.Err)
	}

	if response.Mountpoint != "/etc/docker/mounts/movies" {
		t.Fatal("The mountpoint generated: ", response.Mountpoint, " does not correspond to the proper path")
	}

	mountRequest.Name = "songs"
	mountRequest.ID = "505050"

	response = rdmaVolDriver.Mount(mountRequest)

	if len(response.Err) == 0 {
		t.Fatal("We should not be able to mount and uncreated volume")
	}
}

func TestUnmount(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)
	response := rdmaVolDriver.Unmount(volume.UnmountRequest{Name: "uncreated", ID: "500901234"})

	if len(response.Err) == 0 {
		t.Fatal("We should not be able to Unmount a volume that has not been created yet")
	}

	response = rdmaVolDriver.Unmount(volume.UnmountRequest{})

	if len(response.Err) == 0 {
		t.Fatal("We should not be able to Unmount a volume without a name or ID")
	}

	response = rdmaVolDriver.Create(volume.Request{Name: "Manual Pages"})

	if len(response.Err) != 0 {
		t.Fatal("Encountered error while creating a volume ", response.Err)
	}

	response = rdmaVolDriver.Mount(volume.MountRequest{Name: "Manual Pages", ID: " 909090"})

	if len(response.Err) != 0 {
		t.Fatal("Encountered error while mounting a volume ", response.Err)
	}

	response = rdmaVolDriver.Unmount(volume.UnmountRequest{Name: "Manual Pages", ID: "909090"})

	if len(response.Err) == 0 {
		t.Fatal("There is a difference in IDs when mounting and unmounting, so there should be an error")
	}

	response = rdmaVolDriver.Unmount(volume.UnmountRequest{Name: "Manual Pages", ID: " 909090"})

	if len(response.Err) != 0 {
		t.Fatal(response.Err)
	}
}

func TestCapabilities(t *testing.T) {
	t.Parallel()
	db := db.NewInMemoryVolumeDatabase()
	sc := NewOnDiskStorageController("")

	rdmaVolDriver := NewRDMAVolumeDriver(sc, db)
	response := rdmaVolDriver.Capabilities(volume.Request{})

	if len(response.Err) != 0 {
		t.Fatal(response.Err)
	}
}
