package base

import (
	"flag"
	"testing"
)

var config = flag.String("config", "", "configuration file")

// TestInitializeDatabase: tests the Postmap service
// Requires: a valid database connection and configuration file
// Invoke with: "go test ./... - v -args -config=/valid/path/to/config.json"
func TestInitializeDatabase(t *testing.T) {
	t.Log(*config)
	err := InitializeConfiguration(*config)
	if err != nil {
		t.Error("Error reading configuration file:", err)
	}
	err = InitializeDatabase()
	if err != nil {
		t.Error(err)
	}
	postmapDbs := GetPostmapDb()
	for key, value := range postmapDbs {
		t.Log("Pinging database: ", key, value.Table, value.Key)
		err = value.DbConn.Ping()
		if err != nil {
			t.Error(err)
		}
	}
}

func TestCloseDbConnections(t *testing.T) {
	postmapDb := GetPostmapDb()
	if len(postmapDb) == 0 {
		t.Error("postmapDb should be populated with config file map data")
	}
	closeDbConnections()
	if len(postmapDb) > 0 {
		t.Error("postmapDb should be empty")
	}
}
