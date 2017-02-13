package db

import "testing"

func TestBasic(t *testing.T) {
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
