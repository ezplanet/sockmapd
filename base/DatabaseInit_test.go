package base

import (
	"testing"
)

func TestDbConnections(t *testing.T) {
	err := InitializeConfiguration("../config.json")
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