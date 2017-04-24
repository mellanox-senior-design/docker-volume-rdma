// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/go-plugins-helpers/volume"
)

func TestCreate(t *testing.T) {
	/**************************CREATE********************************/
	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: \"\", Options: map[string]string{}})")
	jsn, err := json.Marshal(volume.Request{
		Name:    "",
		Options: map[string]string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body := bytes.NewBuffer(jsn)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Create", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	var r volume.Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err == "" {
		t.Fatal("Attempting to create a volume with an empty volume_name should cause an error")
	}

	resp.Body.Close()

	/**************************CREATE********************************/
	// Get a Docker approved random name
	volumeName := namesgenerator.GetRandomName(0)

	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s, Options: map[string]string{}})", volumeName)
	jsn, err = json.Marshal(volume.Request{
		Name:    volumeName,
		Options: map[string]string{},
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Create", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to create volume! {Name: ", volumeName, "} ", r.Err)
	}

	resp.Body.Close()

	/**************************CREATE********************************/
	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s, Options: map[string]string{}})", volumeName)
	jsn, err = json.Marshal(volume.Request{
		Name:    volumeName,
		Options: map[string]string{},
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Create", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err == "" {
		t.Fatal("Attempting to create a volume without a unique volume_name should cause an error")
	}

	resp.Body.Close()

	/**************************GET********************************/
	// Get info relating to a paricular volume.
	t.Logf("POST /VolumeDriver.Get")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Get", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to get info for volume: ", volumeName, r.Err)
	}

	if r.Volume.Name != volumeName {
		t.Fatal("Expected: ", volumeName, "Actual:", r.Volume.Name, r.Err)
	}

	resp.Body.Close()

	/**************************REMOVE********************************/
	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create request for server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Remove", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to delete volume! {Name: ", volumeName, "} ", r.Err)
	}
}

func TestGet(t *testing.T) {
	/**************************GET********************************/
	// Get info relating to a paricular volume.
	t.Logf("POST /VolumeDriver.Get")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: \"\", Options: map[string]string{}})")
	jsn, err := json.Marshal(volume.Request{
		Name:    "",
		Options: map[string]string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body := bytes.NewBuffer(jsn)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Get", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	var r volume.Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err == "" {
		t.Fatal("Attempting to get a volume with an empty volume_name should cause an error")
	}

	resp.Body.Close()

	/**************************CREATE********************************/
	// Get a Docker approved random name
	volumeName0 := namesgenerator.GetRandomName(0)

	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s, Options: map[string]string{}})", volumeName0)
	jsn, err = json.Marshal(volume.Request{
		Name:    volumeName0,
		Options: map[string]string{},
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Create", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to create volume! {Name: ", volumeName0, "} ", r.Err)
	}

	resp.Body.Close()

	/**************************GET********************************/
	// Get info relating to a paricular volume.
	t.Logf("POST /VolumeDriver.Get")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName0)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName0,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Get", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to get info for volume: ", volumeName0, r.Err)
	}

	if r.Volume.Name != volumeName0 {
		t.Fatal("Expected: ", volumeName0, "Actual:", r.Volume.Name, r.Err)
	}

	resp.Body.Close()

	/**************************GET********************************/
	// Get a Docker approved random name
	volumeName1 := namesgenerator.GetRandomName(1)

	// Get info relating to a paricular volume.
	t.Logf("POST /VolumeDriver.Get")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName1)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName1,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Get", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err == "" {
		t.Fatal("Volume does not exist so there should be an error, but there is not.")
	}

	/**************************REMOVE********************************/
	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName0)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName0,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Remove", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to delete volume! {Name: ", volumeName0, "} ", r.Err)
	}
}

func TestList(t *testing.T) {
	/**************************LIST********************************/
	// Get the list of volumes registered with the plugin.
	t.Logf("POST /VolumeDriver.List")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{})")
	jsn, err := json.Marshal(volume.Request{})
	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body := bytes.NewBuffer(jsn)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://plugin:8080/VolumeDriver.List", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	var r volume.Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to list volumes!", r.Err)
	}

	if len(r.Volumes) != 0 {
		t.Fatal("List of volumes should be 0! Actual:", len(r.Volumes))
	}

	resp.Body.Close()

	/**************************CREATE********************************/
	// Get a Docker approved random name
	volumeName0 := namesgenerator.GetRandomName(0)

	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s, Options: map[string]string{}})", volumeName0)
	jsn, err = json.Marshal(volume.Request{
		Name:    volumeName0,
		Options: map[string]string{},
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Create", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to create volume! {Name: ", volumeName0, "} ", r.Err)
	}

	resp.Body.Close()

	/**************************LIST********************************/
	// Get the list of volumes registered with the plugin.
	t.Logf("POST /VolumeDriver.List")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{})")
	jsn, err = json.Marshal(volume.Request{})
	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.List", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to list volumes!", r.Err)
	}

	if len(r.Volumes) != 1 {
		t.Fatal("List of volumes should be 1! Actual:", len(r.Volumes))
	}

	resp.Body.Close()

	/**************************CREATE********************************/
	// Get a Docker approved random name
	volumeName1 := namesgenerator.GetRandomName(1)

	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s, Options: map[string]string{}})", volumeName1)
	jsn, err = json.Marshal(volume.Request{
		Name:    volumeName1,
		Options: map[string]string{},
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Create", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to create volume! {Name: ", volumeName1, "} ", r.Err)
	}

	resp.Body.Close()

	/**************************LIST********************************/
	// Get the list of volumes registered with the plugin.
	t.Logf("POST /VolumeDriver.List")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{})")
	jsn, err = json.Marshal(volume.Request{})
	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.List", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to list volumes!", r.Err)
	}

	if len(r.Volumes) != 2 {
		t.Fatal("List of volumes should be 2! Actual:", len(r.Volumes))
	}

	for i := 0; i < len(r.Volumes); i++ {
		name := r.Volumes[i].Name
		if name != volumeName0 && name != volumeName1 {
			t.Fatal("List returned the wrong volume. Expected:", volumeName0, "or", volumeName1, "but received "+r.Volumes[i].Name)
		}
	}

	resp.Body.Close()

	/**************************REMOVE********************************/
	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName0)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName0,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Remove", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to delete volume! {Name: ", volumeName0, "} ", r.Err)
	}

	resp.Body.Close()

	/**************************REMOVE********************************/
	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName1)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName1,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Remove", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to delete volume! {Name: ", volumeName1, "} ", r.Err)
	}
}

func TestRemove(t *testing.T) {
	/**************************REMOVE********************************/
	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: \"\", Options: map[string]string{}})")
	jsn, err := json.Marshal(volume.Request{
		Name:    "",
		Options: map[string]string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body := bytes.NewBuffer(jsn)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Remove", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	var r volume.Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err == "" {
		t.Fatal("Attempting to remove a volume with an empty volume_name should cause an error")
	}

	resp.Body.Close()

	/**************************CREATE********************************/
	// Get a Docker approved random name
	volumeName0 := namesgenerator.GetRandomName(0)

	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s, Options: map[string]string{}})", volumeName0)
	jsn, err = json.Marshal(volume.Request{
		Name:    volumeName0,
		Options: map[string]string{},
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Create", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to create volume! {Name: ", volumeName0, "} ", r.Err)
	}

	resp.Body.Close()

	/**************************CREATE********************************/
	// Get a Docker approved random name
	volumeName1 := namesgenerator.GetRandomName(1)

	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s, Options: map[string]string{}})", volumeName1)
	jsn, err = json.Marshal(volume.Request{
		Name:    volumeName1,
		Options: map[string]string{},
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Create", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to create volume! {Name: ", volumeName1, "} ", r.Err)
	}

	resp.Body.Close()

	/**************************REMOVE********************************/
	// Get a Docker approved random name
	volumeName2 := namesgenerator.GetRandomName(2)

	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName2)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName2,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create request for server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Remove", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err == "" {
		t.Fatal("Error should occur when removing a nonexistant volume")
	}

	resp.Body.Close()

	/**************************REMOVE********************************/
	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName0)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName0,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create request for server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Remove", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to delete volume! {Name: ", volumeName0, "} ", r.Err)
	}

	resp.Body.Close()

	/**************************GET********************************/
	// Get info relating to a paricular volume.
	t.Logf("POST /VolumeDriver.Get")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName0)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName0,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Get", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err == "" {
		t.Fatal("Volume", volumeName0, "should no longer exist", r.Err)
	}

	resp.Body.Close()

	/**************************GET volumeName1********************************/
	// Get info relating to a paricular volume.
	t.Logf("POST /VolumeDriver.Get")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName1)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName1,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Get", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to get info for volume: ", volumeName1, r.Err)
	}

	if r.Volume.Name != volumeName1 {
		t.Fatal("Expected: ", volumeName1, "Actual:", r.Volume.Name, r.Err)
	}

	resp.Body.Close()

	/**************************REMOVE********************************/
	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName1)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName1,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create request for server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Remove", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to delete volume! {Name: ", volumeName1, "} ", r.Err)
	}
}

func TestPath(t *testing.T) {
	// Get a Docker approved random name
	volumeName := namesgenerator.GetRandomName(0)

	/**************************PATH********************************/
	// Path reminds docker of the path a particular volume is attached to.
	t.Logf("POST /VolumeDriver.Path")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName)
	jsn, err := json.Marshal(volume.Request{
		Name: volumeName,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body := bytes.NewBuffer(jsn)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Path", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	var r volume.Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err == "" {
		t.Fatal("Path of a volume that hasn't been created should generate an error")
	}

	resp.Body.Close()

	/**************************CREATE********************************/
	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s, Options: map[string]string{}})", volumeName)
	jsn, err = json.Marshal(volume.Request{
		Name:    volumeName,
		Options: map[string]string{},
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Create", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to create volume! {Name: ", volumeName, "} ", r.Err)
	}

	resp.Body.Close()

	/**************************PATH********************************/
	// Path reminds docker of the path a particular volume is attached to.
	t.Logf("POST /VolumeDriver.Path")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Path", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Error occured while getting path of", volumeName, r.Err)
	}

	if r.Mountpoint != "" {
		t.Fatal("Mountpoint should be: \"\" since volume hasn't been mounted")
	}

	resp.Body.Close()

	/**************************MOUNT********************************/
	// Mount a paricular volume to the host.
	t.Logf("POST /VolumeDriver.Mount")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.MountRequest{Name: %s, ID: \"42\"})", volumeName)
	jsn, err = json.Marshal(volume.MountRequest{
		Name: volumeName,
		ID:   "42",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Mount", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Error occured while mounting:", volumeName, r.Err)
	}

	if r.Mountpoint != "/etc/docker/mounts/"+volumeName {
		t.Fatal("Expected Mountpoint: /etc/docker/mounts/" + volumeName + " Actual: " + r.Mountpoint)
	}

	resp.Body.Close()

	/**************************PATH********************************/
	// Path reminds docker of the path a particular volume is attached to.
	t.Logf("POST /VolumeDriver.Path")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Path", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Error occured while getting path of", volumeName, r.Err)
	}

	if r.Mountpoint != "/etc/docker/mounts/"+volumeName {
		t.Fatal("Expected Mountpoint: /etc/docker/mounts/" + volumeName + " Actual: " + r.Mountpoint)
	}

	resp.Body.Close()

	/**************************UNMOUNT********************************/
	// Unmount a paricular volume from host (if no other mount requests active).
	t.Logf("POST /VolumeDriver.Unmount")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.UnmountRequest{Name: %s, ID: \"42\"})", volumeName)
	jsn, err = json.Marshal(volume.UnmountRequest{
		Name: volumeName,
		ID:   "42",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Unmount", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Error occured while unmounting:", volumeName, r.Err)
	}

	resp.Body.Close()

	/**************************REMOVE********************************/
	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{Name: %s})", volumeName)
	jsn, err = json.Marshal(volume.Request{
		Name: volumeName,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body = bytes.NewBuffer(jsn)
	client = &http.Client{}
	req, err = http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Remove", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch request
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal("Failed to delete volume! {Name: ", volumeName, "} ", r.Err)
	}
}

func TestCapabilities(t *testing.T) {
	// Get the list of capabilities the driver supports.
	t.Logf("POST /VolumeDriver.Capabilities")

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	t.Logf("json.Marshal(volume.Request{})")
	jsn, err := json.Marshal(volume.Request{})
	if err != nil {
		t.Fatal(err)
	}

	// Create request to server
	body := bytes.NewBuffer(jsn)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://plugin:8080/VolumeDriver.Capabilities", body)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch Request
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to connect to server! ", err)
	}

	if resp == nil {
		t.Fatal("resp is nil!")
	}
	if resp.Body == nil {
		t.Fatal("resp.Body is nil!")
	}
	defer resp.Body.Close()

	var r volume.Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Capabilities.Scope != "local" {
		t.Fatal("Scope should be local!")
	}
}
