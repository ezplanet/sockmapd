package service

import (
	"sockmapd/base"
	"sockmapd/model"
	"testing"
)

func TestPostmap(t *testing.T) {
	err := base.InitializeConfiguration("../config.json")
	if err != nil {
		t.Error("Error reading configuration file:", err)
	}
	configuration := base.GetConfiguration()
	if configuration.Database.Port != "3306" {
		t.Error("Could not get correct port number")
	}
	err = base.InitializeDatabase()
	if err != nil {
		t.Error(err)
	}
	var response string
	var request model.Request
	request.Service = "recipient"
	request.Key = "mauro@ezplanet.org"
	response = GetPostmap(request)
	t.Log("Test response: ", response)
	if response != "OK mauro@ezplanet.org" {
		t.Error("Response payload mismatch")
	}
	//
	request.Service = "blacklist"
	request.Key = "hello@buydirect4u.co.uk"
	response = GetPostmap(request)
	t.Log("Test response: ", response)
	if response != "OK REJECT Sender identified as spammer"  {
		t.Error("Response payload mismatch")
	}
	//
	request.Service = "recipient"
	request.Key = "someone@somewhere.com"
	response = GetPostmap(request)
	t.Log("Test response: ", response)
	if response != "NOTFOUND "  {
		t.Error("Response payload mismatch")
	}
	//
	request.Service = "notexists"
	request.Key = "someone@somewhere.com"
	response = GetPostmap(request)
	t.Log("Test response: ", response)
	if response != "TEMP internal configuration error, please try again later"  {
		t.Error("Response payload mismatch")
	}
}