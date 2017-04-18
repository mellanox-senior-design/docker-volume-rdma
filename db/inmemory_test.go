package db

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	t.Parallel()
	im := NewInMemoryVolumeDatabase()

	err := im.Connect()
	if err != nil {
		t.Fatal("Connect returned an error! " + err.Error())
	}

	list, err := im.List()
	if err != nil {
		t.Fatal("List returned an error! " + err.Error())
	}

	if len(list) != 0 {
		t.Fatal("List should have returned exactly zero results.")
	}

	err = im.Create("volumeName", map[string]string{})
	if err != nil {
		t.Fatal("Create returned an error! " + err.Error())
	}

	err = im.Create("volumeName", map[string]string{})
	if err == nil {
		t.Error("Can not create the same volume so we need an error! However, db.Create did not return an error")
	}

	list, err = im.List()
	if err != nil {
		t.Fatal("List returned an error! " + err.Error())
	}

	if len(list) != 1 {
		t.Fatal("List should have returned exactly one result.")
	}

	err = im.Remove("volumeName")
	if err != nil {
		t.Fatal("Remove returned an error! " + err.Error())
	}

	list, err = im.List()
	if err != nil {
		t.Fatal("List returned an error! " + err.Error())
	}

	if len(list) != 0 {
		t.Fatal("List should have returned exactly zero results.")
	}

	err = im.Disconnect()
	if err != nil {
		t.Fatal("Disconnect returned an error! " + err.Error())
	}
}

func TestInMemCreate(t *testing.T) {
	t.Parallel()
	im := NewInMemoryVolumeDatabase()

	err := im.Create("MovieFiles", nil)
	if err != nil {
		t.Fatal("Error obtained while attempting to create volume: ", err)
	}

	err = im.Create("MovieFiles", nil)
	if err == nil {
		t.Error("We should obtain an error when attempting to create a volume that is already created")
	}

	im.Remove("MovieFiles")
}

func TestInMemGet(t *testing.T) {
	t.Parallel()
	im := NewInMemoryVolumeDatabase()

	vol, err := im.Get("Non-existing")
	if err == nil {
		t.Error("Can not get a volume that does not exist")
	}

	err = im.Create("MusicFiles", nil)
	if err != nil {
		t.Fatal("Error obtained while attempting to create volume", err)
	}

	err = im.Create("MovieFiles", nil)
	if err != nil {
		t.Fatal("Error obtained while attempting to create volume", err)
	}

	vol, err = im.Get("MusicFiles")
	if err != nil {
		t.Error(err)
	}

	if vol.Name != "MusicFiles" || len(vol.Mountpoint) != 0 {
		t.Error("Get returned an incorrect volume")
	}

	vol, err = im.Get("MovieFiles")
	if err != nil {
		t.Error(err)
	}

	if vol.Name != "MovieFiles" || len(vol.Mountpoint) != 0 {
		t.Error("Get returned an incorrect volume")
	}

	for i := 0; i <= 100; i++ {
		err = im.Create("VolNum"+strconv.Itoa(i), nil)
		if err != nil {
			t.Fatal("Error obtained while attempting to create volume", i, err)

		}
	}

	for i := 100; i >= 0; i-- {
		volName := "VolNum" + strconv.Itoa(i)
		vol, err = im.Get(volName)

		if vol.Name != volName || len(vol.Mountpoint) != 0 {
			t.Error("Get returned an incorrect volume")
		}
	}
}

func TestInMemPath(t *testing.T) {
	t.Parallel()
	im := NewInMemoryVolumeDatabase()

	err := im.Create("MusicFiles", nil)
	if err != nil {
		t.Fatal("Error obtained while attempting to create volume: ", err)
	}

	err = im.Mount("MusicFiles", "42", "/mnt/sure")
	if err != nil {
		t.Fatal("Error obtained while attempting to mount volume: ", err)
	}

	path, err := im.Path("MusicFiles")
	if err == nil {
		if path != "/mnt/sure" {
			t.Error("Expecting path of mount to be /mnt/sure instead got ", path)
		}
	} else {
		t.Error("Error obtained while attempting to get path of volume: ", err)
	}

	_, err = im.Path("MovieFiles")
	if err == nil {
		t.Error("Expecting error for a volume that is not created")
	}

	err = im.Create("MovieFiles", nil)
	if err != nil {
		t.Fatal("Error obtained while attempting to create volume: ", err)
	}

	path, err = im.Path("MovieFiles")
	if err == nil {
		if len(path) != 0 {
			t.Error("Should not return a path ( ", path, " ) since volume not mounted")
		}
	} else {
		t.Error("Error obtained while attemtping to get path of volume: ", err)
	}
}

func TestInMemRemove(t *testing.T) {
	t.Parallel()
	im := NewInMemoryVolumeDatabase()

	err := im.Remove("Non-existing-Vol")
	if err == nil {
		t.Error("Should not be able to remove a volume that does not exist")
	}

	err = im.Create("MusicFiles", nil)
	if err != nil {
		t.Fatal("Error obtained while attempting to create volume", err)
	}

	err = im.Remove("MusicFiles")
	if err != nil {
		t.Error(err)
	}

	_, err = im.Get("MusicFiles")
	if err == nil {
		t.Error("Volume should not exist, so we need an error")
	}

	err = im.Create("MovieFiles", nil)
	if err != nil {
		t.Fatal("Error obtained while attempting to create volume", err)
	}

	err = im.Mount("MovieFiles", "42", "/mnt/sure/")
	if err != nil {
		t.Fatal("Error obtained while attempting to mount volume", err)
	}

	err = im.Remove("MovieFiles")
	if err == nil {
		t.Error("We should not be able to remove volume if still mounted")
	}

	err = im.Unmount("MovieFiles", "42")
	if err != nil {
		t.Fatal("Error obtained while attempting to unmount volume", err)
	}

	err = im.Remove("MovieFiles")
	if err != nil {
		t.Error(err)
	}

	_, err = im.Get("MovieFiles")
	if err == nil {
		t.Error("Volume should not exist, so we need an error")
	}
}

func TestInMemMount(t *testing.T) {
	t.Parallel()
	im := NewInMemoryVolumeDatabase()

	err := im.Mount("MusicFiles", "43", "/mnt/sure")
	if err == nil {
		t.Error("Should obtain an error if attempting to mount but volume not yet created")
	}

	err = im.Create("MusicFiles", nil)
	if err != nil {
		t.Fatal("Error obtained while attempting to create volume: ", err)
	}

	err = im.Mount("MusicFiles", "42", "/mnt/sure")
	if err != nil {
		t.Fatal("Error obtained while attempting to mount volume: ", err)
	}

	err = im.Mount("MusicFiles", "42", "/mnt/sure/curtain")
	if err != nil {
		t.Fatal("Error obtained while attempting to mount volume: ", err)
	}

	err = im.Mount("MusicFiles", "92", "/mnt/music")
	if err != nil {
		t.Fatal("Error obtained while attempting to mount volume: ", err)
	}

	numMnts := im.mounts["MusicFiles"]["42"]
	assert.Equal(t, 2, numMnts, "MusicFiles with id 42 should be mounted twice")

	numMnts = im.mounts["MusicFiles"]["92"]
	assert.Equal(t, 1, numMnts, "MusicFiles with id 92 should only be mounted once")

}

func TestInMemUnmount(t *testing.T) {
	t.Parallel()
	im := NewInMemoryVolumeDatabase()

	err := im.Unmount("MusicFiles", "42")
	if err == nil {
		t.Error("Should encounter an error when attempting to unmount a volume yet created")
	}

	err = im.Create("MusicFiles", nil)
	if err != nil {
		t.Fatal("Error obtained while attempting to create volume: ", err)
	}

	err = im.Unmount("MusicFiles", "42")
	if err == nil {
		t.Error("Should encounter an error when attempting to unount a volume yet mounted")
	}

	err = im.Mount("MusicFiles", "42", "/mnt/sure")
	if err != nil {
		t.Fatal("Error obtained while attempting to mount volume: ", err)
	}

	err = im.Mount("MusicFiles", "42", "/mnt/sure/curtain")
	if err != nil {
		t.Fatal("Error obtained while attempting to mount volume: ", err)
	}

	err = im.Unmount("MusicFiles", "42")
	if err != nil {
		t.Fatal("Error obtained while unmounting a volume: ", err)
	}

	numMnts := im.mounts["MusicFiles"]["42"]
	assert.Equal(t, 1, numMnts, "Should only have one mount as we previously unmounted")

	err = im.Unmount("MusicFiles", "42")
	if err != nil {
		t.Fatal("Error obtained while unmounting a volume: ", err)
	}

	numMnts = im.mounts["MusicFiles"]["42"]
	assert.Equal(t, 0, numMnts, "Should have no mounts as we previously unmounted")

	err = im.Unmount("MusicFiles", "42")
	if err == nil {
		t.Error("Should encounter an error as volume is no longer mounted")
	}

	err = im.Unmount("MusicFiles", "92")
	if err == nil {
		t.Error("Should encounter an error as volume id not mounted")
	}
}
