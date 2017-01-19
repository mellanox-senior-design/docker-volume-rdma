// Creates a webserver/unix socket that allows a Docker server to create RDMA
// backed volumes.
package main

import (
	"errors"
	"flag"
	"os"
	"strconv"

	"github.com/Jacobingalls/docker-volume-rdma/db"
	"github.com/Jacobingalls/docker-volume-rdma/drivers"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
)

// Port to launch service on.
var httpPort int
var volumeDatabaseDriver string
var volumeDatabasePath string

func init() {
	// Configure application flags.
	flag.IntVar(&httpPort, "port", 8080, "tcp/ip port to serve volume driver on")
	flag.StringVar(&volumeDatabaseDriver, "db", "sqlite", "set the database backend used to store volume metadata: [sqlite, in-memory]")
	flag.StringVar(&volumeDatabasePath, "dbpath", "", "set the database storage path")
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
		if volumeDatabasePath != "" {
			glog.Fatal("Invalid -dbpath, in-memory does not use -dbpath.")
		}
		return db.NewInMemoryVolumeDatabase(), nil

	case "sqlite":
		// If the database path was not set, then we are resposible for managing it.
		if volumeDatabasePath == "" {
			volumeDatabasePath = "./sqlite.db/db"
			glog.Info("Defaulting -dbpath=" + volumeDatabasePath)

			// Create the enclosing dir for the database, if it does not exist.
			_, err := os.Open("./sqlite.db")
			if os.IsNotExist(err) {
				err := os.Mkdir("./sqlite.db", 0755)
				if err != nil {
					glog.Fatal("Unable to create ./sqlite.db folder.")
				}
			}
		}

		return db.NewSqliteVolumeDatabase(volumeDatabasePath), nil

	default:
		return nil, errors.New("Unsupported database, please choose sqlite or in-memory.")
	}
}
