package base

import (
	"log"
	"testing"
)

func TestGetConfiguration(t *testing.T) {
	err := InitializeConfiguration("../config.json")
	if err != nil {
		t.Error("Error reading configuration file:", err)
	}
	conf := GetConfiguration()
	log.Println(conf)
	if conf.Database.Port != "3306" {
		t.Error("Could not get correct port number")
	}
	if len(conf.Postmaps) < 1 {
		t.Error("No postmaps")
	}
	if conf.SysConfig.Port != "2224" {
		t.Error("Socket service port mismatch:", conf.SysConfig.Port)
	}
}