package db

import (
	"database/sql"
	"errors"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
)

// SQLVolumeDatabase defines a volume database that uses a sqlite database.
type SQLVolumeDatabase struct {
	DBType       string
	DBDataSource string
	DBQueries    VolumeDatabaseQueries
	sqlDB        *sql.DB
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

// DefaultSQLQueries stores the default SQL functions for sqldbs to use.
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
// Overide default sql queries by passing a VolumeDatabaseQueries object. Non-nil fields will be be used instead of the defaults.
func NewSQLVolumeDatabase(dbType string, dbDataSource string, dbQueries VolumeDatabaseQueries) SQLVolumeDatabase {

	queries := dbQueries.merge(DefaultSQLQueries)

	return SQLVolumeDatabase{
		DBType:       dbType,
		DBDataSource: dbDataSource,
		DBQueries:    queries,
		sqlDB:        nil}
}

// Connect to database
func (s SQLVolumeDatabase) Connect() error {
	glog.Info("Opening database file: " + s.DBDataSource)

	// Connect to database
	var err error
	s.sqlDB, err = sql.Open(s.DBType, s.DBDataSource)
	if err != nil {
		return err
	}

	err = s.sqlDB.Ping()
	if err != nil {
		glog.Fatal("Unable to connect to database! ", err)
	}

	// Create the volumes table if it do not exist
	glog.Info(s.DBQueries.volumesCreateTableSQL)
	_, err = s.sqlDB.Exec(s.DBQueries.volumesCreateTableSQL)
	if err != nil {
		glog.Error(err, ": ", s.DBQueries.volumesCreateTableSQL)
		return err
	}

	// Create mount table, this will hold all of the ids that are requesting a volume
	glog.Info(s.DBQueries.mountsCreateTableSQL)
	_, err = s.sqlDB.Exec(s.DBQueries.mountsCreateTableSQL)
	if err != nil {
		glog.Error(err, ": ", s.DBQueries.mountsCreateTableSQL)
		return err
	}

	glog.Info("Connected to db.")
	return nil
}

// Disconnect from database
func (s SQLVolumeDatabase) Disconnect() error {
	if err := s.VerifyOrCrash(); err != nil {
		return err
	}

	glog.Info("Closing database: " + s.DBDataSource)
	return s.sqlDB.Close()
}

// VerifyOrCrash if the database connection is not properly configured
func (s SQLVolumeDatabase) VerifyOrCrash() error {
	if s.sqlDB == nil {
		return errors.New("Database is not connected!")
	}

	return nil
}

// Create volume
func (s SQLVolumeDatabase) Create(volumeName string, options map[string]string) error {
	if err := s.VerifyOrCrash(); err != nil {
		return err
	}

	// Verify input.
	if volumeName == "" {
		return errors.New("volume name cannot be empty")
	}

	// Begin transaction to the database
	transaction, err := s.sqlDB.Begin()
	if err != nil {
		return err
	}

	// Prepare the query
	preparedStatement, err := transaction.Prepare(s.DBQueries.volumesInsertSQL)
	if err != nil {
		transaction.Rollback()
		return err
	}
	defer preparedStatement.Close()

	// Actually make the insert
	_, err = preparedStatement.Exec(volumeName)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Commit the change
	return transaction.Commit()
}

// List all volumes
func (s SQLVolumeDatabase) List() ([]*volume.Volume, error) {
	if err := s.VerifyOrCrash(); err != nil {
		return nil, err
	}

	// Query the database about the volumes
	rows, err := s.sqlDB.Query(s.DBQueries.volumesGetNameAndMountpointListSQL)
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
	if err := s.VerifyOrCrash(); err != nil {
		return nil, 0, err
	}

	// Prepare the query
	preparedStatement, err := s.sqlDB.Prepare(s.DBQueries.volumesGetVolumeByNameSQL)
	if err != nil {
		return nil, 0, err
	}
	defer preparedStatement.Close()

	// Query the database about the volumes
	rows, err := preparedStatement.Query(volumeName)
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
		return nil, 0, errors.New("volume does not exist")
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
		return errors.New("volume cannot be removed as it still has active mount requests")
	}

	id, err := s.getVolumeIDByName(volumeName)
	if err != nil {
		return err
	}

	// Begin transaction to the database
	transaction, err := s.sqlDB.Begin()
	if err != nil {
		return err
	}

	// Prepare the queries
	mountsPreparedStatement, err := transaction.Prepare(s.DBQueries.mountsDeleteByVolumeIDSQL)
	if err != nil {
		return err
	}
	defer mountsPreparedStatement.Close()

	volumesPreparedStatement, err := transaction.Prepare(s.DBQueries.volumesDeleteByIDSQL)
	if err != nil {
		return err
	}
	defer volumesPreparedStatement.Close()

	// Actually make the deletes
	_, err = mountsPreparedStatement.Exec(id)
	if err != nil {
		transaction.Rollback()
		return err
	}

	_, err = volumesPreparedStatement.Exec(id)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Commit the change
	return transaction.Commit()
}

// Mount the volume with name and id
func (s SQLVolumeDatabase) Mount(volumeName string, id string, mointpoint string) error {
	mounts, _, err := s.listMounts(volumeName)
	if err != nil {
		return err
	}

	vol, volid, err := s.getVolumeByName(volumeName)
	if err != nil {
		return err
	}

	// Begin transaction to the database
	transaction, err := s.sqlDB.Begin()
	if err != nil {
		return err
	}

	if vol.Mountpoint == "" {
		vol.Mountpoint = mointpoint
		volumeMountPointPreparedStatement, errUp := transaction.Prepare(s.DBQueries.volumesUpdateMountpointSQL)
		if errUp != nil {
			return errUp
		}
		defer volumeMountPointPreparedStatement.Close()

		_, errUp = volumeMountPointPreparedStatement.Exec(vol.Mountpoint, volid)
		if errUp != nil {
			transaction.Rollback()
			return errUp
		}
	}

	var newCount int
	var updateOrInsertQuery string
	numberOfMounts, exists := mounts[id]
	if exists {
		// Need to update the value to number + 1
		newCount = numberOfMounts + 1
		updateOrInsertQuery = s.DBQueries.mountsUpdateCountByVolumeIDAndRequesterSQL
	} else {
		// Need to insert
		newCount = 1
		updateOrInsertQuery = s.DBQueries.mountsInsertSQL
	}

	mountPreparedStatement, err := transaction.Prepare(updateOrInsertQuery)
	if err != nil {
		transaction.Rollback()
		return err
	}
	defer mountPreparedStatement.Close()

	_, err = mountPreparedStatement.Exec(newCount, volid, id)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Commit the change
	return transaction.Commit()
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
	transaction, err := s.sqlDB.Begin()
	if err != nil {
		return err
	}

	var newCount int
	numberOfMounts, exists := mounts[id]
	if exists {
		if numberOfMounts > 1 {

			// Need to update the value to number + 1
			newCount = numberOfMounts - 1
			preparedStatement, err := transaction.Prepare(s.DBQueries.mountsUpdateCountByVolumeIDAndRequesterSQL)
			if err != nil {
				return err
			}
			defer preparedStatement.Close()

			_, err = preparedStatement.Exec(newCount, volid, id)
			if err != nil {
				transaction.Rollback()
				return err
			}

		} else {

			// Need to delete
			preparedStatement, err := transaction.Prepare(s.DBQueries.mountsDeleteByVolumeIDAndRequesterSQL)
			if err != nil {
				return err
			}
			defer preparedStatement.Close()

			_, err = preparedStatement.Exec(volid, id)
			if err != nil {
				transaction.Rollback()
				return err
			}

			volumeMountpointPreparedStatement, err := transaction.Prepare(s.DBQueries.volumesUpdateMountpointSQL)
			if err != nil {
				return err
			}
			defer volumeMountpointPreparedStatement.Close()

			_, err = volumeMountpointPreparedStatement.Exec("", volid)
			if err != nil {
				transaction.Rollback()
				return err
			}
		}
	} else {
		transaction.Rollback()
		return errors.New("volume + ID was not mounted")
	}

	// Commit the change
	return transaction.Commit()
}

// listMounts returns of all the IDs requesting the volume to be mounted and number of requests outstanding for that id.
func (s SQLVolumeDatabase) listMounts(volumeName string) (map[string]int, int, error) {

	id, err := s.getVolumeIDByName(volumeName)
	if err != nil {
		return nil, 0, err
	}

	// Prepare the query
	preparedStatement, err := s.sqlDB.Prepare(s.DBQueries.mountsGetRequesterAndCountByVolumeIDListSQL)
	if err != nil {
		return nil, 0, err
	}
	defer preparedStatement.Close()

	// Query the database about the volumes
	rows, err := preparedStatement.Query(id)
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

func (d VolumeDatabaseQueries) merge(defaults VolumeDatabaseQueries) VolumeDatabaseQueries {

	update := func(overrideValue string, defaultValue string) string {
		if overrideValue != "" {
			return overrideValue
		}

		return defaultValue
	}

	return VolumeDatabaseQueries{
		// Volumes SQL statments
		volumesCreateTableSQL:              update(d.volumesCreateTableSQL, defaults.volumesCreateTableSQL),
		volumesInsertSQL:                   update(d.volumesInsertSQL, defaults.volumesInsertSQL),
		volumesGetNameAndMountpointListSQL: update(d.volumesGetNameAndMountpointListSQL, defaults.volumesGetNameAndMountpointListSQL),
		volumesGetVolumeByNameSQL:          update(d.volumesGetVolumeByNameSQL, defaults.volumesGetVolumeByNameSQL),
		volumesUpdateMountpointSQL:         update(d.volumesUpdateMountpointSQL, defaults.volumesUpdateMountpointSQL),
		volumesDeleteByIDSQL:               update(d.volumesDeleteByIDSQL, defaults.volumesDeleteByIDSQL),

		// Mounts SQL statements
		mountsCreateTableSQL:                        update(d.mountsCreateTableSQL, defaults.mountsCreateTableSQL),
		mountsInsertSQL:                             update(d.mountsInsertSQL, defaults.mountsInsertSQL),
		mountsGetRequesterAndCountByVolumeIDListSQL: update(d.mountsGetRequesterAndCountByVolumeIDListSQL, defaults.mountsGetRequesterAndCountByVolumeIDListSQL),
		mountsUpdateCountByVolumeIDAndRequesterSQL:  update(d.mountsUpdateCountByVolumeIDAndRequesterSQL, defaults.mountsUpdateCountByVolumeIDAndRequesterSQL),
		mountsDeleteByVolumeIDSQL:                   update(d.mountsDeleteByVolumeIDSQL, defaults.mountsDeleteByVolumeIDSQL),
		mountsDeleteByVolumeIDAndRequesterSQL:       update(d.mountsDeleteByVolumeIDAndRequesterSQL, defaults.mountsDeleteByVolumeIDAndRequesterSQL),
	}

}
