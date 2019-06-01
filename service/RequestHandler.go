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
package service

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sockmapd/base"
	"sockmapd/model"
	"strconv"
	"strings"
)

// HandleConnection: takes an incoming tcp connection and processes the request, then encodes the response into a
// netstring and returns the reponse to the calling clients
// request : "[length]:[request],"
// response: "[length]:[response],"
func HandleConnection(conn net.Conn) {
	remote := strings.Split(conn.RemoteAddr().String(), ":")
	netData, err := bufio.NewReader(conn).ReadString(',')
	if err != nil {
		log.Println(err)
		return
	}
	payload := strings.TrimSpace(string(netData))
	//log.Printf("%s - Received: >%s<\n", remote[0], payload)
	var request model.Request
	request, err = parsePostmap(payload)
	if err != nil {
		log.Printf("%s - %s parsing request payload: '%s'", remote[0], base.StrERROR, payload)
		return
	}
	response := GetPostmap(request)
	// Suppress log entry if the Key = KEEPALIVE to avoid flooding log file
	if request.Key != base.StrKEEPALIVE {
		log.Printf("%s - %s:%s - %s", remote[0], request.Service, request.Key, response)
	}
	_, err = conn.Write([]byte(string(encodeResponse(response))))
	if err != nil {
		log.Printf("%s - %s:%s - %s", remote[0], request.Service, request.Key, err)
	}
	_ = conn.Close()
}

// parsePostmap: parses a postmap query string and returns a Request object containing the two elements of the
// query string: service (or map) and key
// request: "[length]:[map] [key],"
func parsePostmap(request string) (model.Request, error) {
	var req model.Request
	request = strings.TrimRight(request, ",")
	if len(request) > 0 {
		payload := strings.SplitN(request, ":", 2)
		size, err := strconv.Atoi(payload[0])
		if err != nil {
			return req, fmt.Errorf("request checksum is not numeric")
		}
		if size != len(payload[1]) {
			return req, fmt.Errorf("request checksum mismatch: %d != %d", size, len(payload[1]))
		}
		elements := strings.Split(payload[1], " ")
		req.Service = elements[0]
		req.Key = elements[1]
		return req, nil
	} else {
		return req, fmt.Errorf("empty request")
	}
}

func encodeResponse(response string) string {
	return fmt.Sprintf("%d:%s,", len(response), response)
}
