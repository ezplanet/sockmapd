package service

import (
	"strconv"
	"strings"
	"testing"
)

func TestRequestHandlerParse (t *testing.T) {
	req, err := parsePostmap ("45:recipient 2321:2321a:11212a:1212b:cc113:::232,")
	if err != nil {
		t.Errorf("Got error: %s", err)
	}

	if req.Service != "recipient" {
		t.Errorf("Service parsed incorrectly")
	}
	if req.Key != "2321:2321a:11212a:1212b:cc113:::232" {
		t.Errorf("Key parsed incorrectly")
	}
	t.Log(req.Service, req.Key)

	req, err = parsePostmap ("41:recipient somebody@somewhere.com,")
	if err != nil {
		t.Logf("Got error: %s", err)
	} else {
		t.Errorf("This should not pass: %s %s", req.Service, req.Key)
	}

	req, err = parsePostmap ("XX:recipient somebody@somedomain.somesuffix,")
	if err != nil {
		t.Logf("Got error: %s", err)
	} else {
		t.Errorf("This should not pass: %s %s", req.Service, req.Key)
	}

	req, err = parsePostmap ("")
	if err != nil {
		t.Logf("Got error: %s", err)
	} else {
		t.Errorf("This should not pass: %s %s", req.Service, req.Key)
	}

	testResponse := "OK somebody@somewhere.com"
	response := strings.TrimRight(encodeResponse(testResponse), ",")
	array := strings.SplitN(response, ":", 2)
	length, err := strconv.Atoi(array[0])
	if err != nil {
		t.Error("Length not numeric")
	}
	if length != len(array[1]) {
		t.Error("Response length mismatch")
	}
	if array[1] != testResponse  {
		t.Error("Response payload mismatch")
	}
}
