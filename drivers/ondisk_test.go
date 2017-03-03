package drivers

import (
	"os"
	"path"
	"testing"
)

func TestNewOnDiskStorageController(t *testing.T) {
	t.Parallel()
	sc := NewOnDiskStorageController("test/controlla")
	if sc.FSPath != "test/controlla" {
		t.Fatal("Storage controller created in wrond dir ", sc.FSPath)
	}
}

func TestSCConnect(t *testing.T) {
	t.Parallel()
	sc := NewOnDiskStorageController("test/controlla")

	err := sc.Connect()

	if err != nil {
		t.Fatal(err)
	}

}

func TestSCDisconnect(t *testing.T) {
	t.Parallel()
	sc := NewOnDiskStorageController("test/controlla")

	err := sc.Disconnect()

	if err != nil {
		t.Fatal(err)
	}

}

func TestSCMount(t *testing.T) {
	t.Parallel()
	sc := NewOnDiskStorageController("test/controlla")
	/*
		mountedPath, err := sc.Mount("")

		if err == nil {
			t.Fatal("A valid mountedPath name should be provided to storage controller")
		}
	*/

	mountedPath, err := sc.Mount("formulavol1")

	if err != nil {
		t.Fatal(err)
	}

	if mountedPath != "test/controlla/formulavol1" {
		t.Fatal("Did not receive expected path, instead got ", mountedPath)
	}

	mountedPath, err = sc.Mount("formulavol1")

	if err != nil {
		t.Fatal(err, " strange, should have been able to mount the same vol")
	}

	if mountedPath != "test/controlla/formulavol1" {
		t.Fatal("If I am mounting the same volume, the same path should be kept ", mountedPath)
	}

	os.Rename(mountedPath, path.Join(path.Dir(mountedPath), path.Base(mountedPath)+".unmounted"))

	mountedPath, err = sc.Mount("formulavol1")

	if err != nil {
		t.Fatal(err, " we should not fail though.")
	}

	if mountedPath != "test/controlla/formulavol1" {
		t.Fatal("Did not expect ", mountedPath, " to be the mountedPath")
	}

}

func TestSCUnmount(t *testing.T) {
	t.Parallel()
	sc := NewOnDiskStorageController("test/controlla")

	err := sc.Unmount("formulavol2")

	if err == nil {
		t.Fatal("There should be an error when unmounting a volume that has not been mounted")
	}

	_, err = sc.Mount("formulavol2")

	if err != nil {
		t.Fatal(err)
	}

	err = sc.Unmount("formulavol2")

	if err != nil {
		t.Fatal(err)
	}

	err = sc.Unmount("formulavol2")

	if err == nil {
		t.Fatal("Should have received an error for unmounting a volume twice")
	}

}

func TestSCDelete(t *testing.T) {
	t.Parallel()
	sc := NewOnDiskStorageController("test/controlla")

	_, err := sc.Mount("volume3")

	if err != nil {
		t.Fatal(err)
	}

	err = sc.Delete("volume3")

	if err != nil {
		t.Fatal(err)
	}

	_, err = sc.Mount("volume4")
	if err != nil {
		t.Fatal(err)
	}
	err = sc.Unmount("volume4")

	err = sc.Delete("volume4")

	if err != nil {
		t.Fatal(err)
	}

	err = sc.Delete("vol5")

	if err != nil {
		t.Fatal(err)
	}
}
