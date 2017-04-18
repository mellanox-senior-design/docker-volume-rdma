package db

import (
	"os"
	"testing"
)

func TestNewSQLiteVolumeDatabase(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		dbPath   string
		expected string
	}{
		{"movies", "movies/db"},
		{"music", "music/db"},
		{"cars", "cars/db"},
		{"houses", "houses/db"},
		{"", "sqlite.db/db"},
	}

	for _, test := range tests {
		actual := NewSQLiteVolumeDatabase(test.dbPath)
		if actual.DBType != "sqlite3" {
			t.Errorf("NewSQLiteVolumeDatabase(%s) = %v; want %v", test.dbPath, actual.DBType, test.expected)
		}

		if actual.DBDataSource != test.expected {
			t.Errorf("NewSQLiteVolumeDatabase(%s) = %v; want %v", test.dbPath, actual.DBDataSource, test.expected)
		}

		if test.dbPath != "" {
			os.RemoveAll(test.dbPath)
		} else {
			os.RemoveAll("sqlite.db")
		}
	}

}

func TestSQLiteConnect(t *testing.T) {
	volDB := NewSQLiteVolumeDatabase("")

	err := volDB.Connect()

	if err != nil {
		t.Error(err)
	}

	os.RemoveAll("sqlite.db")
}
