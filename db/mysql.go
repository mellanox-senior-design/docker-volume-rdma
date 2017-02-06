package db

import (
	"path"

	// Allows connecting to mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

// NewMySQLVolumeDatabase creates a new SQLVolumeDatabase, connectiong to a mysql host.
func NewMySQLVolumeDatabase(host string, username string, password string, schema string) SQLVolumeDatabase {

	var queries VolumeDatabaseQueries
	queries.volumesCreateTableSQL = `CREATE TABLE IF NOT EXISTS volumes (
        id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(256) NOT NULL UNIQUE,
        mountpoint TEXT
    )`

	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	var connection string
	var printedConnection string

	if username == "" {
		username = "root"
	}
	connection = username
	printedConnection = username

	if password != "" {
		connection += ":" + password
		printedConnection += ":<supplied password>"
	}

	connection += "@" + host
	printedConnection += "@" + host

	if schema == "" {
		glog.Fatal("A database schema must be specified with -dbschema")
	}
	connection = path.Join(connection, schema)
	printedConnection = path.Join(printedConnection, schema)

	// Create the connection
	glog.Info("Connecting to", printedConnection)
	return NewSQLVolumeDatabase("mysql", connection, queries)
}
