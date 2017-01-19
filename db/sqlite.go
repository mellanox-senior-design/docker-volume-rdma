package db

import (
	"database/sql"
	"errors"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"

	// Starts sqlite db in the background
	_ "github.com/mattn/go-sqlite3"
)

var sqliteDB *sql.DB

// SqliteVolumeDatabase defines a volume database that uses a sqlite database.
type SqliteVolumeDatabase struct {
	DBPath string
}

// NewSqliteVolumeDatabase creates a new SqliteVolumeDatabase, saving the database at dbPath.
func NewSqliteVolumeDatabase(dbPath string) SqliteVolumeDatabase {
	return SqliteVolumeDatabase{dbPath}
}

// Connect to database
func (s SqliteVolumeDatabase) Connect() error {
	glog.Info("Opening database file: " + s.DBPath)

	// Connect to database
	var err error
	sqliteDB, err = sql.Open("sqlite3", s.DBPath)
	if err != nil {
		return err
	}

	// Create the volumes table if it do not exist
	volumesTableCreateSQL := `CREATE TABLE IF NOT EXISTS volumes (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			mountpoint TEXT
		)`
	_, err = sqliteDB.Exec(volumesTableCreateSQL)
	if err != nil {
		glog.Error(err, ": ", volumesTableCreateSQL)
		return err
	}

	// Create mount table, this will hold all of the ids that are requesting a volume
	mountsTableCreateSQL := `CREATE TABLE IF NOT EXISTS mounts (
			volume_id INTEGER NOT NULL,
			requester_id TEXT NOT NULL,
			count INTEGER NOT NULL
		)`
	_, err = sqliteDB.Exec(mountsTableCreateSQL)
	if err != nil {
		glog.Error(err, ": ", mountsTableCreateSQL)
		return err
	}

	glog.Info("Connected to db.")

	return nil
}

// Disconnect from database
func (s SqliteVolumeDatabase) Disconnect() error {
	return sqliteDB.Close()
}

// Create volume
func (s SqliteVolumeDatabase) Create(volumeName string, options map[string]string) error {

	if sqliteDB == nil {
		glog.Fatal("Database is not connected!")
	}

	if volumeName == "" {
		return errors.New("Volume name cannot be empty.")
	}

	// Begin transaction to the database
	tx, err := sqliteDB.Begin()
	if err != nil {
		return err
	}

	// Prepare the query
	stmt, err := tx.Prepare("INSERT INTO volumes(name) VALUES (?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Actually make the insert
	_, err = stmt.Exec(volumeName)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the change
	return tx.Commit()
}

// List all volumes
func (s SqliteVolumeDatabase) List() ([]*volume.Volume, error) {

	if sqliteDB == nil {
		glog.Fatal("Database is not connected!")
	}

	// Query the database about the volumes
	rows, err := sqliteDB.Query("SELECT name, mountpoint FROM volumes;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Itterate over all of the rows creating the slice of volumes we need.
	var vols []*volume.Volume
	for rows.Next() {
		var name string
		var mountpointNS sql.NullString
		err = rows.Scan(&name, &mountpointNS)
		if err != nil {
			return nil, err
		}

		// If there is a mountpoint, use it.
		var mountpoint string
		if mountpointNS.Valid {
			mountpoint = mountpointNS.String
		}

		// Create volume
		vol := volume.Volume{
			Name:       name,
			Mountpoint: mountpoint,
			Status:     nil}

		// Append it to the list
		vols = append(vols, &vol)
	}

	// Check to see if there was an error durring interation
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return vols, nil
}

// Get information about a particular volue
func (s SqliteVolumeDatabase) Get(volumeName string) (*volume.Volume, error) {
	vol, _, err := s.getVolumeByName(volumeName)
	return vol, err
}

func (s SqliteVolumeDatabase) getVolumeIDByName(volumeName string) (int, error) {
	_, id, err := s.getVolumeByName(volumeName)
	if err == nil {
		glog.Info(volumeName, " is id ", id)
	}
	return id, err
}

func (s SqliteVolumeDatabase) getVolumeByName(volumeName string) (*volume.Volume, int, error) {
	if sqliteDB == nil {
		glog.Fatal("Database is not connected!")
	}

	// Prepare the query
	stmt, err := sqliteDB.Prepare("SELECT id, name, mountpoint FROM volumes WHERE name = ? LIMIT 1;")
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	// Query the database about the volumes
	rows, err := stmt.Query(volumeName)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Itterate over all of the rows creating the slice of volumes we need. (there should only ever be one or zero)
	var vols []*volume.Volume
	var ids []int
	for rows.Next() {
		var id int
		var name string
		var mountpointNS sql.NullString
		err = rows.Scan(&id, &name, &mountpointNS)
		if err != nil {
			return nil, 0, err
		}

		// If there is a mountpoint, use it.
		var mountpoint string
		if mountpointNS.Valid {
			mountpoint = mountpointNS.String
		}

		// Create volume
		vol := volume.Volume{
			Name:       name,
			Mountpoint: mountpoint,
			Status:     nil}

		// Append it to the list
		vols = append(vols, &vol)
		ids = append(ids, id)
	}

	// Check to see if there was an error durring interation
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	// Did we get any results?
	if len(vols) == 0 {
		return nil, 0, errors.New("Volume does not exist.")
	}

	return vols[0], ids[0], nil
}

// Path returns the mountpath of a particular volume
func (s SqliteVolumeDatabase) Path(volumeName string) (string, error) {
	vol, err := s.Get(volumeName)
	if err != nil {
		return "", err
	}

	return vol.Mountpoint, nil
}

// Remove (Delete) a volume from the database.
func (s SqliteVolumeDatabase) Remove(volumeName string) error {
	_, requests, err := s.listMounts(volumeName)
	if err != nil {
		return err
	}

	if requests > 0 {
		return errors.New("Volume cannot be removed as it still has active mount requests.")
	}

	id, err := s.getVolumeIDByName(volumeName)
	if err != nil {
		return err
	}

	// Begin transaction to the database
	tx, err := sqliteDB.Begin()
	if err != nil {
		return err
	}

	// Prepare the queries
	stmtM, err := tx.Prepare("DELETE FROM mounts WHERE volume_id = ?;")
	if err != nil {
		return err
	}
	defer stmtM.Close()

	stmtV, err := tx.Prepare("DELETE FROM volumes WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmtV.Close()

	// Actually make the deletes
	_, err = stmtM.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmtV.Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the change
	return tx.Commit()
}

// Mount the volume with name and id
func (s SqliteVolumeDatabase) Mount(volumeName string, id string) (string, error) {
	mounts, _, err := s.listMounts(volumeName)
	if err != nil {
		return "", err
	}

	vol, volid, err := s.getVolumeByName(volumeName)
	if err != nil {
		return "", err
	}

	// Begin transaction to the database
	tx, err := sqliteDB.Begin()
	if err != nil {
		return "", err
	}

	if vol.Mountpoint == "" {
		vol.Mountpoint = "/etc/docker/mounts/" + volumeName
		stmtUp, errUp := tx.Prepare("UPDATE volumes SET mountpoint=? WHERE id = ?;")
		if errUp != nil {
			return "", errUp
		}
		defer stmtUp.Close()

		_, errUp = stmtUp.Exec(vol.Mountpoint, volid)
		if errUp != nil {
			tx.Rollback()
			return "", errUp
		}
	}

	var newCount int
	var q string
	number, exists := mounts[id]
	if exists {
		// Need to update the value to number + 1
		newCount = number + 1
		q = "UPDATE mounts SET count=? WHERE volume_id = ? and requester_id = ?;"
	} else {
		// Need to insert
		newCount = 1
		q = "INSERT INTO mounts(count, volume_id, requester_id) VALUES (?, ?, ?);"
	}

	stmt, err := tx.Prepare(q)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newCount, volid, id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Commit the change
	return vol.Mountpoint, tx.Commit()
}

func (s SqliteVolumeDatabase) Unmount(volumeName string, id string) error {
	mounts, _, err := s.listMounts(volumeName)
	if err != nil {
		return err
	}

	_, volid, err := s.getVolumeByName(volumeName)
	if err != nil {
		return err
	}

	// Begin transaction to the database
	tx, err := sqliteDB.Begin()
	if err != nil {
		return err
	}

	var newCount int
	number, exists := mounts[id]
	if exists {
		if number > 1 {

			// Need to update the value to number + 1
			newCount = number - 1
			stmt, err := tx.Prepare("UPDATE mounts SET count=? WHERE volume_id = ? and requester_id = ?;")
			if err != nil {
				return err
			}
			defer stmt.Close()

			_, err = stmt.Exec(newCount, volid, id)
			if err != nil {
				tx.Rollback()
				return err
			}

		} else {

			// Need to delete
			stmt, err := tx.Prepare("DELETE FROM mounts WHERE volume_id = ? AND requester_id = ?;")
			if err != nil {
				return err
			}
			defer stmt.Close()

			_, err = stmt.Exec(volid, id)
			if err != nil {
				tx.Rollback()
				return err
			}

			stmtUp, errUp := tx.Prepare("UPDATE volumes SET mountpoint=? WHERE id = ?;")
			if errUp != nil {
				return errUp
			}
			defer stmtUp.Close()

			_, errUp = stmtUp.Exec("", volid)
			if errUp != nil {
				tx.Rollback()
				return errUp
			}
		}
	} else {
		tx.Rollback()
		return errors.New("Volume + ID was not mounted.")
	}

	// Commit the change
	return tx.Commit()
}

// listMounts returns of all the IDs requesting the volume to be mounted and number of requests outstanding for that id.
func (s SqliteVolumeDatabase) listMounts(volumeName string) (map[string]int, int, error) {

	id, err := s.getVolumeIDByName(volumeName)
	if err != nil {
		return nil, 0, err
	}

	// Prepare the query
	stmt, err := sqliteDB.Prepare("SELECT requester_id, count FROM mounts WHERE volume_id = ?;")
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	// Query the database about the volumes
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Itterate over all of the rows creating the slice of volumes we need.
	mounts := map[string]int{}
	sum := 0
	for rows.Next() {
		var requester string
		var count int

		err = rows.Scan(&requester, &count)
		if err != nil {
			return nil, 0, err
		}

		if count > 0 {
			_, exists := mounts[requester]
			if exists {
				mounts[requester] += count
			} else {
				mounts[requester] = count
			}
			sum += count
		}
	}

	// Check to see if there was an error durring interation
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	glog.Info(sum, " mounts for volume ", volumeName)
	return mounts, sum, nil
}
