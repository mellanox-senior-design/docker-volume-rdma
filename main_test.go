package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/mellanox-senior-design/docker-volume-rdma/db"
	"github.com/mellanox-senior-design/docker-volume-rdma/drivers"
)

func configureTest(t *testing.T) (*drivers.RDMAVolumeDriver, *volume.Handler, string, string) {
	// Create temporary directory to "mount" volumes to.
	tempDir, err := ioutil.TempDir("", "docker-volume-rdma")
	if err != nil {
		t.Fatal("Unable to create temp dir! ", err)
	}
	t.Log("Temp Dir: ", tempDir)

	httpPort := strconv.Itoa(10000 + rand.Intn(10000))
	t.Logf("Configuring server for port: %s", httpPort)

	// Configure flags.
	flag.Set("port", httpPort)
	flag.Set("db", "in-memory")
	flag.Set("sc", "on-disk")
	flag.Set("scpath", tempDir)
	flag.Parse()

	// Configure driver and handler.
	configuredDriver, configuredHandler, err := configure()
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that we are using an in memory database, if this fails, check for flag parsing.
	if _, ok := configuredDriver.VolumeDatabase.(db.InMemoryVolumeDatabase); !ok {
		t.Fatal("Configured Driver's VolumeDatabase was not an db.InMemoryVolumeDatabase")
	}

	// Ensure that we are using an on disk storage controller, if this failes, check for flag parsing.
	if _, ok := configuredDriver.StorageController.(drivers.OnDiskStorageController); !ok {
		t.Fatal("Configured Driver's StorageController was not an drivers.OnDiskStorageController")
	}

	if configuredHandler == nil {
		t.Fatal("Configured Handler is nil")
	}

	return configuredDriver, configuredHandler, tempDir, httpPort
}

func TestMain(t *testing.T) {

	_, err := os.Open("/etc/docker")
	if err != nil {
		t.Skip("/etc/docker is required to exist.")
	}

	// Get Configured Driver.
	_, _, tempDir, httpPort := configureTest(t)
	defer os.Remove(tempDir)

	// Start main!
	go main()

	// Create json for request - local variable json masks the global symbol json referring to the JSON module
	jsn, err := json.Marshal(volume.Request{})
	if err != nil {
		t.Fatal(err)
	}

	// Create request to server.
	body := bytes.NewBuffer(jsn)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:"+httpPort+"/VolumeDriver.Capabilities", body)
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

func Test_createListDelete(t *testing.T) {
	// Get Configured Driver.
	driver, _, tempDir, _ := configureTest(t)
	defer os.Remove(tempDir)

	// Get a Docker approved random name
	volumeName := namesgenerator.GetRandomName(0)

	// List all volumes
	t.Logf("driver.List(request.Volume{})")
	response := driver.List(volume.Request{})
	if response.Err != "" {
		t.Fatal("Failed to list volumes! {} ", response.Err)
	}
	if len(response.Volumes) != 0 {
		t.Fatal("driver.List returned ", len(response.Volumes), " volumes. Expected 0.")
	}

	// Create a volume
	t.Logf("driver.Create(volume.Request{Name: %s, Options: map[string]string{}})", volumeName)
	response = driver.Create(volume.Request{
		Name:    volumeName,
		Options: map[string]string{},
	})
	if response.Err != "" { // It is sufficent that there is not an error.
		t.Fatal("Failed to create volume! {Name: ", volumeName, "} ", response.Err)
	}

	// List all volumes
	t.Logf("driver.List(request.Volume{})")
	response = driver.List(volume.Request{})
	if response.Err != "" {
		t.Fatal("Failed to list volumes! {} ", response.Err)
	}
	if len(response.Volumes) != 1 {
		t.Fatal("driver.List returned ", len(response.Volumes), " volumes. Expected 1.")
	}

	// Remove volume
	t.Logf("driver.Remove(volume.Request{Name: %s})", volumeName)
	response = driver.Remove(volume.Request{
		Name: volumeName,
	})
	if response.Err != "" { // It is sufficent that there is not an error.
		t.Fatal("Failed to remove volume! {Name: ", volumeName, "} ", response.Err)
	}

	// List all volumes
	t.Logf("driver.List(request.Volume{})")
	response = driver.List(volume.Request{})
	if response.Err != "" {
		t.Fatal("Failed to list volumes! {} ", response.Err)
	}
	if len(response.Volumes) == 1 {
		t.Fatal("driver.Remove failed to remove volume.")
	} else if len(response.Volumes) != 0 {
		t.Fatal("driver.List returned ", len(response.Volumes), " volumes. Expected 0.")
	}
}

func Test_badConfigurations(t *testing.T) {
	// Test Volume Database
	// Configure flags.
	flag.Set("db", "invalid")
	flag.Set("sc", "on-disk")
	flag.Parse()

	// Configure driver and handler.
	_, _, err := configure()
	if err == nil {
		t.Error("Invalid -db did not cause an error.")
	}

	// Test Storage Controller
	// Configure flags.
	flag.Set("db", "in-memory")
	flag.Set("sc", "invalid")
	flag.Parse()

	// Configure driver and handler.
	_, _, err = configure()
	if err == nil {
		t.Error("Invalid -sc did not cause an error.")
	}
}

func TestGetDatabaseConnection_inmemory(t *testing.T) {
	// Test in-memory Volume Database
	// Configure flags.
	flag.Set("db", "in-memory")
	flag.Set("sc", "on-disk")
	flag.Parse()

	// Configure driver and handler.
	configuredDriver, _, err := configure()
	if err != nil {
		t.Error(err)
	}

	// Ensure that we are using an in memory database, if this fails, check for flag parsing.
	if _, ok := configuredDriver.VolumeDatabase.(db.InMemoryVolumeDatabase); !ok {
		t.Fatal("Configured Driver's VolumeDatabase was not an db.InMemoryVolumeDatabase")
	}
}

func TestGetDatabaseConnection_sqlite(t *testing.T) {
	// Test sqlite Volume Database
	// Configure flags.
	flag.Set("db", "sqlite")
	flag.Set("sc", "on-disk")
	flag.Parse()

	// Configure driver and handler.
	configuredDriver, _, err := configure()
	if err != nil {
		t.Error(err)
	}

	// Ensure that we are using an in memory database, if this fails, check for flag parsing.
	if _, ok := configuredDriver.VolumeDatabase.(db.SQLVolumeDatabase); !ok {
		t.Fatal("Configured Driver's VolumeDatabase was not an db.SQLVolumeDatabase")
	}
}

func TestGetDatabaseConnection_mysql(t *testing.T) {
	// Test mysql Volume Database
	// Configure flags.
	flag.Set("db", "mysql")
	flag.Set("dbschema", "rdma")
	flag.Set("sc", "on-disk")
	flag.Parse()

	// Configure driver and handler.
	configuredDriver, _, err := configure()
	if err != nil {
		t.Error(err)
	}

	// Ensure that we are using an in memory database, if this fails, check for flag parsing.
	if _, ok := configuredDriver.VolumeDatabase.(db.SQLVolumeDatabase); !ok {
		t.Fatal("Configured Driver's VolumeDatabase was not an db.SQLVolumeDatabase")
	}
}

func TestGetStorageConnection_ondisk(t *testing.T) {
	// Test OnDiskStorageController Volume Database
	// Configure flags.
	flag.Set("sc", "on-disk")
	flag.Parse()

	// Configure driver and handler.
	configuredDriver, _, err := configure()
	if err != nil {
		t.Error(err)
	}

	// Ensure that we are using an in memory database, if this fails, check for flag parsing.
	if _, ok := configuredDriver.StorageController.(drivers.OnDiskStorageController); !ok {
		t.Fatal("Configured Driver's StorageController was not an drivers.OnDiskStorageController")
	}
}

func TestGetStorageConnection_glusterfs(t *testing.T) {
	// Test GlusterStorageController Volume Database
	// Configure flags.
	flag.Set("sc", "glusterfs")
	flag.Parse()

	// Configure driver and handler.
	configuredDriver, _, err := configure()
	if err != nil {
		t.Error(err)
	}

	// Ensure that we are using an in memory database, if this fails, check for flag parsing.
	if _, ok := configuredDriver.StorageController.(drivers.GlusterStorageController); !ok {
		t.Fatal("Configured Driver's StorageController was not an drivers.GlusterStorageController")
	}
}
