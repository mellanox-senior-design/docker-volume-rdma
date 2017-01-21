// Creates a webserver/unix socket that allows a Docker server to create RDMA
// backed volumes.
package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Jacobingalls/docker-volume-rdma/db"
	"github.com/Jacobingalls/docker-volume-rdma/drivers"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
)

// Port to launch service on.
var httpPort int
var volumeDatabaseDriver string
var volumeDatabasePath string
var volumeDatabaseHost string
var volumeDatabaseUsername string
var volumeDatabasePassword string
var volumeDatabaseSchema string

func init() {
	// Configure application flags.
	flag.IntVar(&httpPort, "port", 8080, "tcp/ip port to serve volume driver on")

	flag.StringVar(&volumeDatabaseDriver, "db", "sqlite", "set the database backend used to store volume metadata: [sqlite, mysql, in-memory]")
	flag.StringVar(&volumeDatabasePath, "dbpath", "", "set the database storage path")

	flag.StringVar(&volumeDatabaseHost, "dbhost", "", "set the database host (default is '' or localhost:3306)")
	flag.StringVar(&volumeDatabaseUsername, "dbuser", "", "set the database username (default is root)")
	flag.StringVar(&volumeDatabasePassword, "dbpass", "", "set the database password (optional)")
	flag.StringVar(&volumeDatabaseSchema, "dbschema", "", "set the database schema (required)")
}

// Configure and start the docker volume plugin server.
func main() {

	// Parse flags as glog needs the flags to be solidified before starting.
	flag.Parse()

	// Convert port to string, and print startup message.
	port := strconv.Itoa(httpPort)

	// Create and begin serving volume driver on tcp/ip port, httpPort.
	vd, dberr := GetDatabaseConnection()
	if dberr != nil {
		glog.Fatal(dberr)
	}

	// Configure Storage Controller
	sc := drivers.NewGlusterStorageController()

	// Print startup message and start server
	glog.Info("Connecting to services ...")
	driver := drivers.NewRDMAVolumeDriver(sc, vd)
	driver.Connect()
	defer driver.Disconnect()

	// Handle SIGINT gracefully
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		glog.Info("Exiting ...")
		driver.Disconnect()
		os.Exit(1)
	}()

	glog.Info("Running! http://localhost:" + port)
	h := volume.NewHandler(driver)
	err := h.ServeTCP("test_volume", ":"+port, nil)

	// Log any error that may have occurred.
	glog.Fatal(err)
}

// GetDatabaseConnection returns the database connection that was requested on the command line.
func GetDatabaseConnection() (db.VolumeDatabase, error) {
	switch volumeDatabaseDriver {
	case "in-memory":
		validateDatabaseFlags(true, true, false, false, false, false)
		return db.NewInMemoryVolumeDatabase(), nil

	case "sqlite":
		validateDatabaseFlags(false, true, false, false, false, false)
		return db.NewSQLiteVolumeDatabase(volumeDatabasePath), nil

	case "mysql":
		validateDatabaseFlags(false, false, true, true, true, true)
		return db.NewMySQLVolumeDatabase(volumeDatabaseHost, volumeDatabaseUsername, volumeDatabasePassword, volumeDatabaseSchema), nil

	default:
		return nil, errors.New("Unsupported database, please choose sqlite, mysql, or in-memory.")
	}
}

func validateDatabaseFlags(fatal bool, path bool, host bool, username bool, password bool, schema bool) {

	var errors bool
	noteErrorFunc := func(name string, value string, used bool) {
		if !used && value != "" {
			glog.Warning("Volume Driver: ", volumeDatabaseDriver, " does not support ", name, ".")
			errors = true
		}
	}

	noteErrorFunc("-dbpath", volumeDatabasePath, path)
	noteErrorFunc("-dbhost", volumeDatabaseHost, host)
	noteErrorFunc("-dbuser", volumeDatabaseUsername, username)
	noteErrorFunc("-dbpass", volumeDatabasePassword, password)
	noteErrorFunc("-dbschema", volumeDatabaseSchema, schema)

	if errors && fatal {
		glog.Fatal("Invalid flag(s) were passed, are you using the correct volume driver?")
	}
}
