// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
)

func TestMain(t *testing.T) {
	// Get the list of capabilities the driver supports.
	t.Logf("POST /VolumeDriver.Capabilities")

	// Create json for request
	// note: a local variable json masks the global symbol json referring to the JSON module
	// renamed json to jsn for now
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

	body = new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	jsn = body.Bytes()
	err = json.Unmarshal(jsn, &r)
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
	// note: a local variable json masks the global symbol json referring to the JSON module
	// renamed json to jsn for now
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

	body = new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	jsn = body.Bytes()
	err = json.Unmarshal(jsn, &r)
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

func TestListCreateListAndDelete(t *testing.T) {
}
