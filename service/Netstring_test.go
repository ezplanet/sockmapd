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

//
// This class includes the test cases for netstring encoder and decoder
//

package service

import (
	"sockmapd/base"
	"strconv"
	"strings"
	"testing"
)

const MyString = "my string 28 characters long"
const MyStringEncoded = "28:my string 28 characters long,"

const MyRequestEncoded = "45:recipient 2321:2321a:11212a:1212b:cc113:::232,"
const MyRequestLengthMismatch = "99:recipient 2321:2321a:11212a:1212b:cc113:::232,"
const MyRequestLengthNotNum = "XX:recipient 2321:2321a:11212a:1212b:cc113:::232,"
const MyRequest = "recipient 2321:2321a:11212a:1212b:cc113:::232"

// TestNetstringEncode: tests the service to encode a string into a netstring
func TestNetStringEncode(t *testing.T) {
	netstring := NetStringEncode(MyString)
	if string(netstring) != MyStringEncoded {
		t.Errorf("Netstring mismatch: '%s' != '%s'", string(netstring), MyStringEncoded)
	}

	testResponse := "OK somebody@somewhere.com"
	netstring = NetStringEncode(testResponse)
	response := strings.TrimRight(string(netstring), ",")
	array := strings.SplitN(response, ":", 2)
	length, err := strconv.Atoi(array[0])
	if err != nil {
		t.Error("Length not numeric")
	}
	if length != len(array[1]) {
		t.Error("Response length mismatch")
	}
	if array[1] != testResponse {
		t.Error("Response payload mismatch")
	}
}

func TestNetStringDecode(t *testing.T) {
	netstring, err := NetStringDecode(MyRequestEncoded)
	if err != nil {
		t.Errorf("Got error: %s", err)
	}
	if netstring != MyRequest {
		t.Errorf("Decoded string mismatch: '%s' != '%s'", netstring, MyRequest)
	}

	netstring, err = NetStringDecode(MyRequestLengthNotNum)
	if err == nil || err.Error() != base.ErrNsLenNotNumeric {
		t.Errorf("Should report %s", base.ErrNsLenNotNumeric)
	}

	netstring, err = NetStringDecode(MyRequestLengthMismatch)
	if err == nil || err.Error() != base.ErrNsLenMismatch+": 99 != 45" {
		t.Errorf("Should report %s != %s", base.ErrNsLenMismatch, err.Error())
	}

	netstring, err = NetStringDecode("")
	if err == nil || err.Error() != base.ErrNsEmpty {
		t.Errorf("Should report %s", base.ErrNsEmpty)
	}

	netstring, err = NetStringDecode("  ,  ")
	if err == nil || err.Error() != base.ErrNsInvalid {
		t.Errorf("Should report %s", base.ErrNsInvalid)
	}
}
