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
package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sockmapd/base"
	"sockmapd/service"
	"strconv"
)

func main() {
	var port string
	var config string
	flag.StringVar(&port, "p", base.TcpPORT, "Please specify a valid port number between 1025 and 65535")
	flag.StringVar(&config, "c", "config.json", "Please specify a valid configuration file path")
	flag.Parse()

	tcpPort, err := strconv.Atoi(port)
	if err != nil || tcpPort > 65534 || tcpPort < 1025 {
		log.Fatalf("TCP Port '%s' is not a number between 1025 and 65535", port)
	}
	if _, err := os.Stat(config); os.IsNotExist(err) {
		log.Fatalf("Configuration file '%s' is not readable or is not present", config)
	}

	err = base.InitializeConfiguration(config)
	if err != nil {
		log.Fatal(err, base.StrTerminated)
	}
	conf := base.GetConfiguration()
	if len(conf.SysConfig.Logfile) > 0 {
		f, err := os.OpenFile(conf.SysConfig.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			log.Printf("%s opening log file: %v", base.StrERROR, err)
		} else {
			defer f.Close()
			log.SetOutput(f)
		}
	}
	if len(conf.SysConfig.Port) > 0 {
		port = conf.SysConfig.Port
	}
	err = base.InitializeDatabase()
	if err != nil {
		log.Fatalln(err, base.StrTerminated)
	}
	base.ProcessSignal()
	log.Println("Listening to port: ", port)
	listener, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("%s: %s", base.StrERROR, err)
			return
		}
		service.HandleConnection(conn)
	}
}
