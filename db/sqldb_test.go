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
