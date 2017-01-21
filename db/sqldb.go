package db

import (
	"database/sql"
	"errors"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
)

var sqlDB *sql.DB

// SQLVolumeDatabase defines a volume database that uses a sqlite database.
type SQLVolumeDatabase struct {
	DBType       string
	DBDataSource string
	DBQueries    VolumeDatabaseQueries
}

// VolumeDatabaseQueries is a struct that can contain queries for the volume database
type VolumeDatabaseQueries struct {
	// Volumes SQL statements
	volumesCreateTableSQL              string
	volumesInsertSQL                   string
	volumesGetNameAndMountpointListSQL string
	volumesGetVolumeByNameSQL          string
	volumesUpdateMountpointSQL         string
	volumesDeleteByIDSQL               string

	// Mounts SQL statements
	mountsCreateTableSQL                        string
	mountsInsertSQL                             string
	mountsGetRequesterAndCountByVolumeIDListSQL string
	mountsUpdateCountByVolumeIDAndRequesterSQL  string
	mountsDeleteByVolumeIDSQL                   string
	mountsDeleteByVolumeIDAndRequesterSQL       string
}

var DefaultSQLQueries = VolumeDatabaseQueries{
	// Volumes SQL statments
	volumesCreateTableSQL: `CREATE TABLE IF NOT EXISTS volumes (
        id INTEGER NOT NULL PRIMARY KEY,
        name VARCHAR(256) NOT NULL UNIQUE,
        mountpoint TEXT
    );`,
	volumesInsertSQL:                   "INSERT INTO volumes(name) VALUES (?);",
	volumesGetNameAndMountpointListSQL: "SELECT name, mountpoint FROM volumes;",
	volumesGetVolumeByNameSQL:          "SELECT id, name, mountpoint FROM volumes WHERE name = ? LIMIT 1;",
	volumesUpdateMountpointSQL:         "UPDATE volumes SET mountpoint=? WHERE id = ?;",
	volumesDeleteByIDSQL:               "DELETE FROM volumes WHERE id = ?;",

	// Mounts SQL statements
	mountsCreateTableSQL: `CREATE TABLE IF NOT EXISTS mounts (
        volume_id INTEGER NOT NULL,
        requester_id VARCHAR(256) NOT NULL,
        count INTEGER NOT NULL
    );`,
	mountsInsertSQL:                             "INSERT INTO mounts(count, volume_id, requester_id) VALUES (?, ?, ?);",
	mountsGetRequesterAndCountByVolumeIDListSQL: "SELECT requester_id, count FROM mounts WHERE volume_id = ?;",
	mountsUpdateCountByVolumeIDAndRequesterSQL:  "UPDATE mounts SET count=? WHERE volume_id = ? and requester_id = ?;",
	mountsDeleteByVolumeIDSQL:                   "DELETE FROM mounts WHERE volume_id = ?;",
	mountsDeleteByVolumeIDAndRequesterSQL:       "DELETE FROM mounts WHERE volume_id = ? AND requester_id = ?;",
}

// NewSQLVolumeDatabase creates a new SQLVolumeDatabase, saving the database at dbPath.
func NewSQLVolumeDatabase(dbType string, dbDataSource string, dbQueries VolumeDatabaseQueries) SQLVolumeDatabase {

	queries := dbQueries.merge(DefaultSQLQueries)

	return SQLVolumeDatabase{
		DBType:       dbType,
		DBDataSource: dbDataSource,
		DBQueries:    queries}
}

// Connect to database
func (s SQLVolumeDatabase) Connect() error {
	glog.Info("Opening database file: " + s.DBDataSource)

	// Connect to database
	var err error
	sqlDB, err = sql.Open(s.DBType, s.DBDataSource)
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		glog.Fatal("Unable to connect to database! ", err)
	}

	// Create the volumes table if it do not exist
	_, err = sqlDB.Exec(s.DBQueries.volumesCreateTableSQL)
	if err != nil {
		glog.Error(err, ": ", s.DBQueries.volumesCreateTableSQL)
		return err
	}

	// Create mount table, this will hold all of the ids that are requesting a volume
	_, err = sqlDB.Exec(s.DBQueries.mountsCreateTableSQL)
	if err != nil {
		glog.Error(err, ": ", s.DBQueries.mountsCreateTableSQL)
		return err
	}

	glog.Info("Connected to db.")
	return nil
}

// Disconnect from database
func (s SQLVolumeDatabase) Disconnect() error {
	glog.Info("Closing database file: " + s.DBDataSource)
	return sqlDB.Close()
}

// VerifyOrCrash if the database connection is not properly configured
func (s SQLVolumeDatabase) VerifyOrCrash() {
	if sqlDB == nil {
		glog.Fatal("Database is not connected!")
	}
}

// Create volume
func (s SQLVolumeDatabase) Create(volumeName string, options map[string]string) error {
	s.VerifyOrCrash()

	// Verify input.
	if volumeName == "" {
		return errors.New("Volume name cannot be empty.")
	}

	// Begin transaction to the database
	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}

	// Prepare the query
	stmt, err := tx.Prepare(s.DBQueries.volumesInsertSQL)
	if err != nil {
		tx.Rollback()
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
func (s SQLVolumeDatabase) List() ([]*volume.Volume, error) {
	s.VerifyOrCrash()

	// Query the database about the volumes
	rows, err := sqlDB.Query(s.DBQueries.volumesGetNameAndMountpointListSQL)
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
func (s SQLVolumeDatabase) Get(volumeName string) (*volume.Volume, error) {
	vol, _, err := s.getVolumeByName(volumeName)
	return vol, err
}

func (s SQLVolumeDatabase) getVolumeIDByName(volumeName string) (int, error) {
	_, id, err := s.getVolumeByName(volumeName)
	if err == nil {
		glog.Info(volumeName, " is id ", id)
	}
	return id, err
}

func (s SQLVolumeDatabase) getVolumeByName(volumeName string) (*volume.Volume, int, error) {
	s.VerifyOrCrash()

	// Prepare the query
	stmt, err := sqlDB.Prepare(s.DBQueries.volumesGetVolumeByNameSQL)
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
func (s SQLVolumeDatabase) Path(volumeName string) (string, error) {
	vol, err := s.Get(volumeName)
	if err != nil {
		return "", err
	}

	return vol.Mountpoint, nil
}

// Remove (Delete) a volume from the database.
func (s SQLVolumeDatabase) Remove(volumeName string) error {
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
	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}

	// Prepare the queries
	stmtM, err := tx.Prepare(s.DBQueries.mountsDeleteByVolumeIDSQL)
	if err != nil {
		return err
	}
	defer stmtM.Close()

	stmtV, err := tx.Prepare(s.DBQueries.volumesDeleteByIDSQL)
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
func (s SQLVolumeDatabase) Mount(volumeName string, id string) (string, error) {
	mounts, _, err := s.listMounts(volumeName)
	if err != nil {
		return "", err
	}

	vol, volid, err := s.getVolumeByName(volumeName)
	if err != nil {
		return "", err
	}

	// Begin transaction to the database
	tx, err := sqlDB.Begin()
	if err != nil {
		return "", err
	}

	if vol.Mountpoint == "" {
		vol.Mountpoint = "/etc/docker/mounts/" + volumeName
		stmtUp, errUp := tx.Prepare(s.DBQueries.volumesUpdateMountpointSQL)
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
		q = s.DBQueries.mountsUpdateCountByVolumeIDAndRequesterSQL
	} else {
		// Need to insert
		newCount = 1
		q = s.DBQueries.mountsInsertSQL
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

// Unmount volume with name and id
func (s SQLVolumeDatabase) Unmount(volumeName string, id string) error {
	mounts, _, err := s.listMounts(volumeName)
	if err != nil {
		return err
	}

	_, volid, err := s.getVolumeByName(volumeName)
	if err != nil {
		return err
	}

	// Begin transaction to the database
	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}

	var newCount int
	number, exists := mounts[id]
	if exists {
		if number > 1 {

			// Need to update the value to number + 1
			newCount = number - 1
			stmt, err := tx.Prepare(s.DBQueries.mountsUpdateCountByVolumeIDAndRequesterSQL)
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
			stmt, err := tx.Prepare(s.DBQueries.mountsDeleteByVolumeIDAndRequesterSQL)
			if err != nil {
				return err
			}
			defer stmt.Close()

			_, err = stmt.Exec(volid, id)
			if err != nil {
				tx.Rollback()
				return err
			}

			stmtUp, errUp := tx.Prepare(s.DBQueries.volumesUpdateMountpointSQL)
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
func (s SQLVolumeDatabase) listMounts(volumeName string) (map[string]int, int, error) {

	id, err := s.getVolumeIDByName(volumeName)
	if err != nil {
		return nil, 0, err
	}

	// Prepare the query
	stmt, err := sqlDB.Prepare(s.DBQueries.mountsGetRequesterAndCountByVolumeIDListSQL)
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

func (dest VolumeDatabaseQueries) merge(src VolumeDatabaseQueries) VolumeDatabaseQueries {

	update := func(dst string, src string) string {
		if dst != "" {
			return dst
		}

		return src
	}

	return VolumeDatabaseQueries{
		// Volumes SQL statments
		volumesCreateTableSQL:              update(dest.volumesCreateTableSQL, src.volumesCreateTableSQL),
		volumesInsertSQL:                   update(dest.volumesInsertSQL, src.volumesInsertSQL),
		volumesGetNameAndMountpointListSQL: update(dest.volumesGetNameAndMountpointListSQL, src.volumesGetNameAndMountpointListSQL),
		volumesGetVolumeByNameSQL:          update(dest.volumesGetVolumeByNameSQL, src.volumesGetVolumeByNameSQL),
		volumesUpdateMountpointSQL:         update(dest.volumesUpdateMountpointSQL, src.volumesUpdateMountpointSQL),
		volumesDeleteByIDSQL:               update(dest.volumesDeleteByIDSQL, src.volumesDeleteByIDSQL),

		// Mounts SQL statements
		mountsCreateTableSQL:                        update(dest.mountsCreateTableSQL, src.mountsCreateTableSQL),
		mountsInsertSQL:                             update(dest.mountsInsertSQL, src.mountsInsertSQL),
		mountsGetRequesterAndCountByVolumeIDListSQL: update(dest.mountsGetRequesterAndCountByVolumeIDListSQL, src.mountsGetRequesterAndCountByVolumeIDListSQL),
		mountsUpdateCountByVolumeIDAndRequesterSQL:  update(dest.mountsUpdateCountByVolumeIDAndRequesterSQL, src.mountsUpdateCountByVolumeIDAndRequesterSQL),
		mountsDeleteByVolumeIDSQL:                   update(dest.mountsDeleteByVolumeIDSQL, src.mountsDeleteByVolumeIDSQL),
		mountsDeleteByVolumeIDAndRequesterSQL:       update(dest.mountsDeleteByVolumeIDAndRequesterSQL, src.mountsDeleteByVolumeIDAndRequesterSQL),
	}

}
