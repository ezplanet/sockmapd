// Copyright 2019 by mauro@ezplanet.org (Mauro Mozzarelli)
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package service

import (
	"testing"
)

func TestRequestHandlerParse(t *testing.T) {
	req, err := parsePostmap("45:recipient 2321:2321a:11212a:1212b:cc113:::232,")
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

	req, err = parsePostmap("41:recipient somebody@somewhere.com,")
	if err != nil {
		t.Logf("Got error: %s", err)
	} else {
		t.Errorf("This should not pass: %s %s", req.Service, req.Key)
	}

	req, err = parsePostmap("22:somebody@somewhere.com,")
	if err != nil {
		t.Logf("Got error: %s", err)
		if err.Error() != "invalid request format" {
			t.Errorf("Error should be: 'invalid request format'")
			t.Fail()
		}
	} else {
		t.Errorf("This should not pass: %s %s", req.Service, req.Key)
	}

	req, err = parsePostmap("XX:recipient somebody@somedomain.somesuffix,")
	if err != nil {
		t.Logf("Got error: %s", err)
	} else {
		t.Errorf("This should not pass: %s %s", req.Service, req.Key)
	}

	req, err = parsePostmap("")
	if err != nil {
		t.Logf("Got error: %s", err)
	} else {
		t.Errorf("This should not pass: %s %s", req.Service, req.Key)
	}
}
