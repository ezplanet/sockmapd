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
package main

import (
	"flag"
	"log"
	"math"
	"net"
	"os"
	"os/user"
	"sockmapd/base"
	"sockmapd/logger"
	"sockmapd/service"
	"strconv"
	"syscall"
)

func main() {
	var port string
	var config string
	var logfile string
	var loglevel string

	flag.StringVar(&port, "p", base.TcpPORT, "Please specify a valid port number between 1025 and 65535")
	flag.StringVar(&config, "c", "config.json", "Please specify a valid configuration file path")
	flag.StringVar(&logfile, "l", base.LogFile, "Please specify a valid log file path or stdout/stderr")
	flag.StringVar(&loglevel, "v", base.LogLevel, "Please specify log level: trace, debug, info, warn, error, fatal, panic")
	flag.Parse()

	err := logger.Init(loglevel, logfile)
	if err != nil {
		panic(err)
	}

	tcpPort, err := strconv.Atoi(port)
	if err != nil || tcpPort > 65534 || tcpPort < 1025 {
		logger.Log.Fatal().Msgf("TCP Port '%s' is not a number between 1025 and 65535", port)
	}
	if _, err := os.Stat(config); os.IsNotExist(err) {
		logger.Log.Fatal().Msgf("Configuration file '%s' is not readable or is not present", config)
	}

	err = base.InitializeConfiguration(config)
	if err != nil {
		log.Fatal(err, base.StrTerminated)
	}
	conf := base.GetConfiguration()

	// Dropping user and group privileges (group must be first)
	if len(conf.SysConfig.RunAsGroup) > 0 {
		groupStrut, err := user.LookupGroup(conf.SysConfig.RunAsGroup)
		if err != nil {
			logger.Log.Fatal().Msgf("Error %s, could not run as group '%s' - terminating", err, conf.SysConfig.RunAsGroup)
		} else {
			gid, err := strconv.Atoi(groupStrut.Gid)
			if err != nil {
				logger.Log.Fatal().Msgf("Error %s, could not run as group '%s' - terminating", err, conf.SysConfig.RunAsGroup)
			}
			err = syscall.Setgid(gid)
			if err != nil {
				logger.Log.Fatal().Msgf("Error %s, could not set the gid - terminating", err)
			}
		}
	}

	if len(conf.SysConfig.RunAsUser) > 0 {
		userStrut, err := user.Lookup(conf.SysConfig.RunAsUser)
		if err != nil {
			logger.Log.Fatal().Msgf("Error %s, could not run as user '%s' - terminating", err, conf.SysConfig.RunAsUser)
		} else {
			uid, err := strconv.Atoi(userStrut.Uid)
			if err != nil {
				logger.Log.Fatal().Msgf("Error %s, could not run as user '%s' - terminating", err, conf.SysConfig.RunAsUser)
			}
			err = syscall.Setuid(uid)
			if err != nil {
				logger.Log.Fatal().Msgf("Error %s, could not set the uid - terminating", err)
			}
		}
	}

	if len(conf.SysConfig.Port) > 0 {
		port = conf.SysConfig.Port
	}
	err = base.InitializeDatabase()
	if err != nil {
		logger.Log.Fatal().Msgf("%s, %s", err, base.StrTerminated)
	}
	base.ProcessSignal()
	logger.Log.Info().Msgf("Listening to port: %s", port)
	listener, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		logger.Log.Fatal().Msgf("Could not listen to tcp port %s, %s", port, err)
	}
	var count int = 0
	defer listener.Close()
	for {
		logger.Log.Trace().Msgf("Listening to port: %s for new TCP connection", port)
		conn, err := listener.Accept()
		if err != nil {
			logger.Log.Fatal().Msgf("TCP connection error: %s", err)
		}
		logger.Log.Trace().Msgf("New TCP connection from %s", conn.RemoteAddr())
		if count >= math.MaxInt {
			count = 0
		}
		count++
		go service.HandleConnection(conn, count)
	}
}
