package db

import "testing"

func TestNewMySQLVolumeDatabase(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		host     string
		username string
		password string
		schema   string
		expected string
	}{
		{"localhost", "root", "pass", "schema", "root:pass@localhost/schema"},
		{"localhost", "", "pass", "schema", "root:pass@localhost/schema"},
		{"", "root", "pass", "schema", "root:pass@/schema"},
		{"localhost", "root", "pass", "", ""},
	}
	for _, test := range tests {
		actual, err := NewMySQLVolumeDatabase(test.host, test.username, test.password, test.schema)
		if test.expected == "" {
			if err == nil {
				t.Errorf("NewMySQLVolumeDatabase(%s,%s,%s,%s) = nil; want error", test.host, test.username, test.password, test.schema)
			}
		} else {
			if actual.DBType != "mysql" {
				t.Errorf("NewMySQLVolumeDatabase(%s,%s,%s,%s) = %v; want %v", test.host, test.username, test.password, test.schema, actual.DBType, test.expected)
			}

			if actual.DBDataSource != test.expected {
				t.Errorf("NewMySQLVolumeDatabase(%s,%s,%s,%s) = %v; want %v", test.host, test.username, test.password, test.schema, actual.DBDataSource, test.expected)
			}
		}
	}
}
