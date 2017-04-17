package db

import (
	"database/sql"
	"errors"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewSQLVolumeDatabase_default(t *testing.T) {
	t.Parallel()

	var queries VolumeDatabaseQueries

	foo := NewSQLVolumeDatabase("type", "datasource", queries)
	if foo.DBType != "type" {
		t.Error("type was not configured")
	}

	if foo.DBDataSource != "datasource" {
		t.Error("datasource was not configured")
	}

	if foo.DBQueries != DefaultSQLQueries {
		t.Error("queries are not set to the defaults")
	}
}

func TestNewSQLVolumeDatabase_modifications(t *testing.T) {
	t.Parallel()

	queries := VolumeDatabaseQueries{
		// Volumes SQL statements
		volumesCreateTableSQL:              "a",
		volumesInsertSQL:                   "b",
		volumesGetNameAndMountpointListSQL: "c",
		volumesGetVolumeByNameSQL:          "d",
		volumesUpdateMountpointSQL:         "e",
		volumesDeleteByIDSQL:               "f",

		// Mounts SQL statements
		mountsCreateTableSQL:                        "g",
		mountsInsertSQL:                             "h",
		mountsGetRequesterAndCountByVolumeIDListSQL: "i",
		mountsUpdateCountByVolumeIDAndRequesterSQL:  "j",
		mountsDeleteByVolumeIDSQL:                   "k",
		mountsDeleteByVolumeIDAndRequesterSQL:       "l",
	}

	foo := NewSQLVolumeDatabase("type", "datasource", queries)
	if foo.DBType != "type" {
		t.Error("type was not configured")
	}

	if foo.DBDataSource != "datasource" {
		t.Error("datasource was not configured")
	}

	if foo.DBQueries == DefaultSQLQueries {
		t.Error("queries are set to the defaults")
	}

	if foo.DBQueries != queries {
		t.Error("queries are not set to the overrides")
	}
}

func createMockVolumeDatabase(t *testing.T) (*sql.DB, sqlmock.Sqlmock, VolumeDatabase) {
	// Create mock for a random sql db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// now we execute our method
	volumeDatabase := NewSQLVolumeDatabase("mock", "mock", VolumeDatabaseQueries{})
	sqlDB = db

	return db, mock, volumeDatabase
}

func TestValidateOrCrash(t *testing.T) {
	volumeDatabase := NewSQLVolumeDatabase("mock", "mock", VolumeDatabaseQueries{})

	var tests = []struct {
		name string
		f    func() error
	}{
		{"Disconnect", func() error {
			return volumeDatabase.Disconnect()
		}},

		{"Create", func() error {
			return volumeDatabase.Create("volumeName", map[string]string{})
		}},

		{"List", func() error {
			_, err := volumeDatabase.List()
			return err
		}},

		{"Get", func() error {
			_, err := volumeDatabase.Get("volumeName")
			return err
		}},

		{"Path", func() error {
			_, err := volumeDatabase.Path("volumeName")
			return err
		}},

		{"Remove", func() error {
			return volumeDatabase.Remove("volumeName")
		}},

		{"Mount", func() error {
			return volumeDatabase.Mount("volumeName", "id", "mointpoint")
		}},

		{"Unmount", func() error {
			return volumeDatabase.Unmount("volumeName", "id")
		}},
	}
	for _, test := range tests {
		if err := test.f(); err == nil {
			t.Errorf("%s should have returned an error.", test.name)
		}
	}
}

func TestCreate_trivial(t *testing.T) {
	createSQL := `[INSERT INTO volumes(name) VALUES (?);]`

	db, mock, volumeDatabase := createMockVolumeDatabase(t)
	defer db.Close()

	// Configure Mock
	mock.ExpectBegin()
	prepared := mock.ExpectPrepare(createSQL)
	prepared.ExpectExec().WithArgs("volume_name").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := volumeDatabase.Create("volume_name", map[string]string{}); err != nil {
		t.Errorf("error was not expected while creating volume: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestCreate_badNoName(t *testing.T) {
	t.Parallel()

	db, _, volumeDatabase := createMockVolumeDatabase(t)
	defer db.Close()

	err := volumeDatabase.Create("", map[string]string{})
	if err == nil {
		t.Error("error was expected while creating volume as name cannot be nil.")
	}

	if err.Error() != "volume name cannot be empty" {
		t.Error("error msg was expected to be 'volume name cannot be empty', was '" + err.Error() + "'")
	}
}

func TestCreate_failBegin(t *testing.T) {
	t.Parallel()

	db, mock, volumeDatabase := createMockVolumeDatabase(t)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(errors.New("ExampleError"))

	err := volumeDatabase.Create("volume_name", map[string]string{})
	if err == nil {
		t.Error("error was expected while creating volume as name cannot be nil.")
	}

	if err.Error() != "ExampleError" {
		t.Error("error msg was expected to be 'ExampleError', was '" + err.Error() + "'")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestCreate_failPrepare(t *testing.T) {
	createSQL := `[INSERT INTO volumes(name) VALUES (?);]`

	db, mock, volumeDatabase := createMockVolumeDatabase(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(createSQL).WillReturnError(errors.New("ExampleError"))
	mock.ExpectRollback()

	err := volumeDatabase.Create("volume_name", map[string]string{})
	if err == nil {
		t.Error("error was expected while creating volume as name cannot be nil.")
	}

	if err.Error() != "ExampleError" {
		t.Error("error msg was expected to be 'ExampleError', was '" + err.Error() + "'")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestCreate_failExec(t *testing.T) {
	createSQL := `[INSERT INTO volumes(name) VALUES (?);]`

	db, mock, volumeDatabase := createMockVolumeDatabase(t)
	defer db.Close()

	mock.ExpectBegin()
	prepare := mock.ExpectPrepare(createSQL)
	prepare.ExpectExec().WithArgs("volume_name").WillReturnError(errors.New("ExampleError"))
	mock.ExpectRollback()

	err := volumeDatabase.Create("volume_name", map[string]string{})
	if err == nil {
		t.Error("error was expected while creating volume as name cannot be nil.")
	}

	if err.Error() != "ExampleError" {
		t.Error("error msg was expected to be 'ExampleError', was '" + err.Error() + "'")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestSQLDisconnect(t *testing.T) {
	_, mock, volDB := createMockVolumeDatabase(t)

	mock.ExpectClose()

	err := volDB.Disconnect()

	if err != nil {
		t.Error("was not expecting error at closure, but here it is: ", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestSQLList(t *testing.T) {
	db, mock, volDB := createMockVolumeDatabase(t)

	defer db.Close()

	rows := sqlmock.NewRows([]string{"name", "mountpoint"})
	query := "SELECT name, mountpoint FROM volumes;"

	mock.ExpectQuery(query).WillReturnError(errors.New("db err, can't perform"))
	volList, err := volDB.List()
	if err == nil {
		t.Error("was expecting error from DB")
	} else if err.Error() != "db err, can't perform" {
		t.Error("was not expecting this err thrown: ", err)
	}

	mock.ExpectQuery(query).WillReturnRows(rows)

	volList, err = volDB.List()

	if err != nil {
		t.Error("Not expecting error during listing of empty list of volumes: ", err)
	}

	if len(volList) != 0 {
		t.Error("there should be no volumes returned")
	}

	newRows := sqlmock.NewRows([]string{"name", "mountpoint"}).
		AddRow("aventura_vol", "/etc/mnt/").
		AddRow("movies_vol", "/etc/mnt/").
		AddRow("computer_vol", "")

	mock.ExpectQuery(query).WillReturnRows(newRows)

	volList, err = volDB.List()

	if err != nil {
		t.Error("Not expecting error during listing of volumes: ", err)
	}

	if len(volList) != 3 {
		t.Error("expected to List 3 volumes instead listed ", len(volList))
	}

	for i := 0; i < len(volList); i++ {
		switch volList[i].Name {
		case "aventura_vol", "movies_vol":
			if volList[i].Mountpoint != "/etc/mnt/" {
				t.Error("for ", volList[i].Name, " the mntpoint was incorrectly returned as ", volList[i].Mountpoint)
			}
		case "computer_vol":
			if volList[i].Mountpoint != "" {
				t.Error("for ", volList[i].Name, " the mntpoint was incorrectly returned as ", volList[i].Mountpoint)
			}
		default:
			t.Error("Was not expecting to list a volume named ", volList[i].Name)
		}
	}

	// we make sure that all expectations were met
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestSQLGet(t *testing.T) {
	db, mock, volDB := createMockVolumeDatabase(t)

	defer db.Close()
	query := `[SELECT id, name, mountpoint FROM volumes WHERE name = aventura_vol LIMIT 1;]`

	prepare := mock.ExpectPrepare(query)
	prepare.ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"id", "name", "mountpoint"}))

	_, err := volDB.Get("aventura_vol")
	if err == nil {
		t.Error("was expecting error from 'Get' as there are no volumes")
	}

	prepare = mock.ExpectPrepare(query)
	rows := sqlmock.NewRows([]string{"id", "name", "mountpoint"}).
		AddRow("50", "aventura_vol", "/etc/mnt/")

	prepare.ExpectQuery().WillReturnRows(rows)

	vol, err := volDB.Get("aventura_vol")
	if err != nil {
		t.Error(err)
	}

	if vol.Name != "aventura_vol" {
		t.Error("Did not expect to 'Get' volume with name ", vol.Name)
	}

	mock.ExpectPrepare(query).WillReturnError(errors.New("preperation error"))

	_, err = volDB.Get("aventura_vol")

	if err.Error() != "preperation error" {
		t.Error("expected 'Get' to catch specific error thrown. instead got ", err)
	}

	prepare = mock.ExpectPrepare(query)
	prepare.ExpectQuery().WillReturnError(errors.New("query err"))

	_, err = volDB.Get("aventura_vol")

	if err.Error() != "query err" {
		t.Error("expected 'Get' to catch specific error thrown. instead got ", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

}

func TestSQLPath(t *testing.T) {
	db, mock, volDB := createMockVolumeDatabase(t)

	defer db.Close()
	query := `[SELECT id, name, mountpoint FROM volumes WHERE name = aventura_vol LIMIT 1;]`

	prepare := mock.ExpectPrepare(query)
	rows := sqlmock.NewRows([]string{"id", "name", "mountpoint"}).
		AddRow("42", "aventura_vol", "")

	prepare.ExpectQuery().WillReturnRows(rows)

	path, err := volDB.Path("aventura_vol")

	if err != nil {
		t.Error("Not expecting error while getting Path ", err)
	}

	if path != "" {
		t.Error(path, " was not the path we were expecting")
	}

	prepare = mock.ExpectPrepare(query)
	newRow := sqlmock.NewRows([]string{"id", "name", "mountpoint"}).
		AddRow("42", "aventura_vol", "/etc/mnt/")

	prepare.ExpectQuery().WillReturnRows(newRow)

	path, err = volDB.Path("aventura_vol")

	if err != nil {
		t.Error("Not expecting error while getting Path ", err)
	}

	if path != "/etc/mnt/" {
		t.Error(path, " was not the path we were expecting")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

type responseRows struct {
	id         string
	name       string
	mountpoint string
	requester  string
	count      int
}

func handleGetVolumeByName(mock sqlmock.Sqlmock, prepareException, prepareStatementException bool, prepareErr, prepareSErr string, rowsNeeded []responseRows) {
	query := `[SELECT id, name, mountpoint FROM volumes WHERE name = ? LIMIT 1;]`
	prepare := mock.ExpectPrepare(query)

	if prepareException {
		prepare.WillReturnError(errors.New(prepareErr))
	} else {
		if prepareStatementException {
			prepare.ExpectQuery().WillReturnError(errors.New(prepareSErr))
		} else {
			rows := sqlmock.NewRows([]string{"id", "name", "mountpoint"})
			for i := 0; i < len(rowsNeeded); i++ {
				rows = rows.AddRow(rowsNeeded[i].id, rowsNeeded[i].name, rowsNeeded[i].mountpoint)
			}

			prepare.ExpectQuery().WillReturnRows(rows)
		}
	}
}

func handleListMounts(mock sqlmock.Sqlmock, prepareErr, prepareSErr bool, perr, pserr string, rowsNeeded []responseRows) {
	query := `[SELECT requester_id, count FROM mount WHERE volume_id = ?;]`

	// so list mounts calls getVolumeIDByName
	// which in turn calls getVolumeByName

	handleGetVolumeByName(mock, false, false, "", "", rowsNeeded)

	prepare := mock.ExpectPrepare(query)

	if prepareErr {
		prepare.WillReturnError(errors.New(perr))
	} else {
		if prepareSErr {
			prepare.ExpectQuery().WillReturnError(errors.New(pserr))
		} else {
			rows := sqlmock.NewRows([]string{"requester", "count"})

			for i := 0; i < len(rowsNeeded); i++ {
				rows = rows.AddRow(rowsNeeded[i].requester, rowsNeeded[i].count)
			}

			prepare.ExpectQuery().WillReturnRows(rows)
		}
	}
}

func TestSQLRemove(t *testing.T) {
	db, mock, volDB := createMockVolumeDatabase(t)

	defer db.Close()

	rRows := []responseRows{
		{id: "42", name: "aventura_vol", mountpoint: "/etc/mnt/", requester: "42", count: 5},
	}

	handleListMounts(mock, false, false, "", "", rRows)

	err := volDB.Remove("aventura_vol")

	if err == nil {
		t.Error("we should get an error for attempting to remove volume that is still mounted somewhere")
	} else if err.Error() != "volume cannot be removed as it still has active mount requests" {
		t.Error("we got a total different error than expected : ", err)
	}

	rRows = []responseRows{
		{id: "42", name: "aventura_vol"},
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, true, false, "preperation error", "", rRows)

	err = volDB.Remove("aventura_vol")

	if err == nil {
		t.Error("we should have gotten an error from database")
	} else if err.Error() != "preperation error" {
		t.Error("did not receive the expected error, instead : ", err)
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	mock.ExpectBegin().WillReturnError(errors.New("Database Begin Error"))

	err = volDB.Remove("aventura_vol")

	if err == nil {
		t.Error("we should have gotten an error from database")
	} else if err.Error() != "Database Begin Error" {
		t.Error("did not receive the expected error, instead : ", err)
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	mountQuery := `[DELETE FROM mounts WHERE volume_id = ?;]`

	mock.ExpectBegin()
	mountPrep := mock.ExpectPrepare(mountQuery).WillReturnError(errors.New("prep err"))

	err = volDB.Remove("aventura_vol")

	if err == nil {
		t.Error("we should have gotten an error from database")
	} else if err.Error() != "prep err" {
		t.Error("did not receive the expected error, instead : ", err)
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	volQuery := `[DELETE FROM volumes WHERE id = ?;]`

	mock.ExpectBegin()
	mountPrep = mock.ExpectPrepare(mountQuery)
	volPrep := mock.ExpectPrepare(volQuery).WillReturnError(errors.New("volprep err"))

	err = volDB.Remove("aventura_vol")

	if err == nil {
		t.Error("we should have gotten an error from database")
	} else if err.Error() != "volprep err" {
		t.Error("did not receive the expected error, instead : ", err)
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	mock.ExpectBegin()
	mountPrep = mock.ExpectPrepare(mountQuery)
	volPrep = mock.ExpectPrepare(volQuery)
	mountPrep.ExpectExec().WithArgs(42).WillReturnError(errors.New("mnt delete err"))
	mock.ExpectRollback()

	err = volDB.Remove("aventura_vol")

	if err == nil {
		t.Error("we should have gotten an error from database")
	} else if err.Error() != "mnt delete err" {
		t.Error("did not receive the expected error, instead : ", err)
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	mock.ExpectBegin()
	mountPrep = mock.ExpectPrepare(mountQuery)
	volPrep = mock.ExpectPrepare(volQuery)
	mountPrep.ExpectExec().WithArgs(42).WillReturnResult(sqlmock.NewResult(1, 1))
	volPrep.ExpectExec().WithArgs(42).WillReturnError(errors.New("vol delete err"))
	mock.ExpectRollback()

	err = volDB.Remove("aventura_vol")

	if err == nil {
		t.Error("we should have gotten an error from database")
	} else if err.Error() != "vol delete err" {
		t.Error("did not receive the expected error, instead : ", err)
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	mock.ExpectBegin()
	mountPrep = mock.ExpectPrepare(mountQuery)
	volPrep = mock.ExpectPrepare(volQuery)
	mountPrep.ExpectExec().WithArgs(42).WillReturnResult(sqlmock.NewResult(1, 1))
	volPrep.ExpectExec().WithArgs(42).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = volDB.Remove("aventura_vol")

	if err != nil {
		t.Error("received an error but not expecting one : ", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestSQLMount(t *testing.T) {
	db, mock, volDB := createMockVolumeDatabase(t)

	defer db.Close()

	rRows := []responseRows{
		{id: "42", name: "aventura_vol"},
	}

	handleListMounts(mock, false, false, "", "", rRows)

	handleGetVolumeByName(mock, false, true, "", "error thrown in db!", rRows)

	err := volDB.Mount("aventura_vol", "42", "/etc/mnt/")

	if err == nil {
		t.Error("expecting an error being thrown from DB")
	} else if err.Error() != "error thrown in db!" {
		t.Error("error thrown not expected : ", err)
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	mock.ExpectBegin().WillReturnError(errors.New("db begin error"))
	err = volDB.Mount("aventura_vol", "42", "/etc/mnt/")

	if err == nil {
		t.Error("expecting an error being thrown from DB")
	} else if err.Error() != "db begin error" {
		t.Error("error thrown not expected : ", err)
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	updateMntPt := `[UPDATE volumes SET mounpoint=? WHERE id = ?;]`
	mock.ExpectBegin()
	uMntPrepare := mock.ExpectPrepare(updateMntPt)
	uMntPrepare.ExpectExec().WithArgs("/etc/mnt", 42).WillReturnResult(sqlmock.NewResult(1, 1))

	insMnt := `[INSERT INTO mounts(count, volume_id, requester_id) VALUES (?, ?, ?);]`
	insertPrep := mock.ExpectPrepare(insMnt)
	insertPrep.ExpectExec().WithArgs(1, 42, "42").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = volDB.Mount("aventura_vol", "42", "/etc/mnt")
	if err != nil {
		t.Error("was not expecting an error while mounting, but here it is: ", err)
	}

	rRows = []responseRows{
		{id: "42", name: "aventura_vol", requester: "42", mountpoint: "/etc/mnt", count: 1},
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	updateMntPt = `[UPDATE volumes SET mounpoint=? WHERE id = ?;]`
	mock.ExpectBegin()
	updateCnt := `[UPDATE mounts SET count = ? WHERE volume_id = ? and requester_id = ?;]`
	updateCntPrep := mock.ExpectPrepare(updateCnt)
	updateCntPrep.ExpectExec().WithArgs(2, 42, "42").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = volDB.Mount("aventura_vol", "42", "/etc/mnt")

	if err != nil {
		t.Error(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

}

func TestSQLUnmount(t *testing.T) {
	db, mock, volDB := createMockVolumeDatabase(t)
	defer db.Close()

	rRows := []responseRows{
		{id: "42", name: "aventura_vol", mountpoint: "/etc/mnt/", requester: "42", count: 1},
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	mock.ExpectBegin()

	delete := `[DELETE FROM mounts WHERE volume_id = ? AND requester_id = ?;]`
	delPrepare := mock.ExpectPrepare(delete)
	delPrepare.ExpectExec().WithArgs(42, "42").WillReturnResult(sqlmock.NewResult(1, 1))

	updateMntPt := `[UPDATE volumes SET mountpoint=? WHERE id = ?;]`
	uPrep := mock.ExpectPrepare(updateMntPt)
	uPrep.ExpectExec().WithArgs("", 42).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := volDB.Unmount("aventura_vol", "42")

	if err != nil {
		t.Error("error encountered while unmounting: ", err)
	}

	rRows = []responseRows{
		{id: "50", name: "car_vol", mountpoint: "/etc/mnt", requester: "50", count: 10},
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	mock.ExpectBegin()

	upCnt := `[UPDATE mounts SET count=? WHERE volume_id = ? and requester_id = ?;]`
	upCntPrepare := mock.ExpectPrepare(upCnt)
	upCntPrepare.ExpectExec().WithArgs(9, 50, "50").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = volDB.Unmount("car_vol", "50")

	if err != nil {
		t.Error("error encountered while unmounting: ", err)
	}

	rRows = []responseRows{
		{id: "9", name: "music_vol"},
	}

	handleListMounts(mock, false, false, "", "", rRows)
	handleGetVolumeByName(mock, false, false, "", "", rRows)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = volDB.Unmount("music_vol", "9")

	if err == nil {
		t.Error("we should get an error when attemtping to unmount a volume that is not mounted")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
