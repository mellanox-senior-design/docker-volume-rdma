package db

import (
	"os"
	"path"

	"github.com/golang/glog"
	// Starts sqlite db in the background
	_ "github.com/mattn/go-sqlite3"
)

// NewSQLiteVolumeDatabase creates a new SQLVolumeDatabase, saving the database at dbPath.
func NewSQLiteVolumeDatabase(dbPath string) SQLVolumeDatabase {

	var queries VolumeDatabaseQueries
	queries.volumesCreateTableSQL = `CREATE TABLE IF NOT EXISTS volumes (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        mountpoint TEXT
    )`

	queries.mountsCreateTableSQL = `CREATE TABLE IF NOT EXISTS mounts (
        volume_id INTEGER NOT NULL,
        requester_id TEXT NOT NULL,
        count INTEGER NOT NULL
    )`

	// If the database path was not set, then we are resposible for managing it.
	if dbPath == "" {
		dbPath = "./sqlite.db"
		glog.Info("Defaulting -dbpath=" + dbPath)
	}

	// Create the enclosing dir for the database, if it does not exist.
	_, err := os.Open(dbPath)
	if os.IsNotExist(err) {
		glog.Info("Creating: ", dbPath)
		err := os.MkdirAll(dbPath, 0755)
		if err != nil {
			glog.Fatal("Unable to create ", dbPath, " folder.")
		}
	}

	// Create the connection
	return NewSQLVolumeDatabase("sqlite3", path.Join(dbPath, "db"), VolumeDatabaseQueries{})
}
