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
	"bufio"
	"fmt"
	"log"
	"net"
	"sockmapd/base"
	"sockmapd/model"
	"strings"
)

// HandleConnection: takes an incoming tcp connection and processes the request, then encodes the response into a
// netstring and returns the reponse to the calling clients
// request : "[length]:[request],"
// response: "[length]:[response],"
func HandleConnection(conn net.Conn) {
	var response string
	remote := strings.Split(conn.RemoteAddr().String(), ":")
	for {
		netData, err := bufio.NewReader(conn).ReadString(',')
		if err != nil {
			if err.Error() != "EOF" {
				log.Println(err)
			}
			break
		}
		var request model.Request
		request, err = parsePostmap(netData)
		if err != nil {
			log.Printf("%s - %s: %s - parsing request: '%s'", remote[0], base.StrERROR, err, request)
			response = base.SmapPERM + err.Error()
		} else {
			response = GetPostmap(request)
			// Suppress log entry if the Key = KEEPALIVE to avoid flooding log file
			if request.Key != base.StrKEEPALIVE {
				log.Printf("%s - %s:%s - %s", remote[0], request.Service, request.Key, response)
			}
		}
		_, err = conn.Write([]byte(NetStringEncode(response)))
		if err != nil {
			log.Printf("%s - %s:%s - %s", remote[0], request.Service, request.Key, err)
		}
	}
	_ = conn.Close()
}

// parsePostmap: parses a postmap query string and returns a Request object containing the two elements of the
// query string: service (or map) and key
// request: "[length]:[map] [key],"
func parsePostmap(request string) (model.Request, error) {
	var req model.Request
	payload, err := NetStringDecode(request)
	if err != nil {
		return req, err
	}
	elements := strings.Split(payload, " ")
	if len(elements) != 2 {
		return req, fmt.Errorf(base.ErrParseInvalid)
	}
	req.Service = elements[0]
	req.Key = elements[1]
	return req, nil
}
