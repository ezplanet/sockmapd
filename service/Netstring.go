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
// This class includes the services necessary to handle the encoding and decoding of netstrings
//

package service

import (
	"fmt"
	"sockmapd/base"
	"strconv"
	"strings"
)

// Encode: takes a regular string as input and returns a netstring https://en.wikipedia.org/wiki/Netstring
// that is a byte encoded string with the format "[len]:[mystring],"
// input : "my string"
// return: [length]:[mystring],"
func NetStringEncode(input string) string {
	netstring := fmt.Sprintf("%d:%s,", len(input), input)
	return netstring
}

// NetStringDecode: decodes and validates a netstring, returns the payload or an error
func NetStringDecode(netstring string) (string, error) {
	netstring = strings.TrimRight(netstring, ",")
	if len(netstring) > 0 {
		payload := strings.SplitN(netstring, ":", 2)
		if len(payload) != 2 {
			return "", fmt.Errorf(base.ErrNsInvalid)
		}
		size, err := strconv.Atoi(payload[0])
		if err != nil {
			return "", fmt.Errorf(base.ErrNsLenNotNumeric)
		}
		if size != len(payload[1]) {
			return "", fmt.Errorf("%s: %d != %d", base.ErrNsLenMismatch, size, len(payload[1]))
		}
		return payload[1], nil
	} else {
		return "", fmt.Errorf(base.ErrNsEmpty)
	}
}
