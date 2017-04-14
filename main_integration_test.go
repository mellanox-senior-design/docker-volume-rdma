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

func TestCapabilities(t *testing.T) {
	// Get the list of capabilities the driver supports.
	t.Logf("POST /VolumeDriver.Capabilities")

	// Create json for request
	// local variable json masks the global symbol json referring to the JSON module
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
	if resp == nil || resp.Body == nil {
		t.Fatal("nil response!")
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

func TestList(t *testing.T) {
	// Get the list of volumes registered with the plugin.
	t.Logf("POST /VolumeDriver.List")

	// Create json for request
	// local variable json masks the global symbol json referring to the JSON module
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
	if resp == nil || resp.Body == nil {
		t.Fatal("nil response!")
	}
	defer resp.Body.Close()

	var r volume.Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal(r.Err)
	}
	if len(r.Volumes) != 0 {
		t.Fatal("List of volumes should be 0!")
	}
}

func TestListCreateListRemoveList(t *testing.T) {
	/**************************LIST********************************/
	// Get the list of volumes registered with the plugin.
	t.Logf("POST /VolumeDriver.List")

	// Create json for request
	// local variable json masks the global symbol json referring to the JSON module
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
	if resp == nil || resp.Body == nil {
		t.Fatal("nil response!")
	}

	var r volume.Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal(r.Err)
	}
	if len(r.Volumes) != 0 {
		t.Fatal("List of volumes should be 0!")
	}
	resp.Body.Close()

	/**************************CREATE********************************/
	// Create a new volume with name and options.
	t.Logf("POST /VolumeDriver.Create")

	// Get a Docker approved random name
	volumeName := namesgenerator.GetRandomName(0)

	// Create json for request
	// local variable json masks the global symbol json referring to the JSON module
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
	if resp == nil || resp.Body == nil {
		t.Fatal("nil response!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal(r.Err)
	}
	resp.Body.Close()

	/**************************LIST********************************/
	// Get the list of volumes registered with the plugin.
	t.Logf("POST /VolumeDriver.List")

	// Create json for request
	// local variable json masks the global symbol json referring to the JSON module
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
	if resp == nil || resp.Body == nil {
		t.Fatal("nil response!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal(r.Err)
	}
	if len(r.Volumes) != 1 {
		t.Fatal("List of volumes should be 1!")
	}
	resp.Body.Close()

	/**************************REMOVE********************************/
	// Remove (Delete) a paricular volume.
	t.Logf("POST /VolumeDriver.Remove")

	// Create json for request
	// local variable json masks the global symbol json referring to the JSON module
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
	if resp == nil || resp.Body == nil {
		t.Fatal("nil response!")
	}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal(r.Err)
	}
	resp.Body.Close()

	/**************************LIST********************************/
	// Get the list of volumes registered with the plugin.
	t.Logf("POST /VolumeDriver.List")

	// Create json for request
	// local variable json masks the global symbol json referring to the JSON module
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
	if resp == nil || resp.Body == nil {
		t.Fatal("nil response!")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.Err != "" {
		t.Fatal(r.Err)
	}
	if len(r.Volumes) != 0 {
		t.Fatal("List of volumes should be 0!")
	}
}
