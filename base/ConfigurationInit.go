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
package base

import (
	"encoding/json"
	"fmt"
	"os"
	"sockmapd/model"
	"strconv"
)

var configuration model.Configuration
var configFile string

// InitializeConfiguration:
// configurationFile is an  optional parameter. It is called. This function is called with the
// configuration file as parameter by main() at startup, subsequently the file path is saved and
// re-used when the function is called again when a SIGHUP or SIGUSR1 is triggered.
// The configuration file is opened and the content copied into the Configuration object.
func InitializeConfiguration(configurationFile ...string) error {
	if len(configurationFile) == 1 {
		configFile = configurationFile[0]
	}
	file, _ := os.Open(configFile)
	defer file.Close()
	decoder := json.NewDecoder(file)
	var conf model.Configuration
	err := decoder.Decode(&conf)
	if err != nil {
		return fmt.Errorf("%s decoding configuration file: %s", StrERROR, err)
	} else {
		if len(conf.SysConfig.Port) > 0 {
			tcpPort, err := strconv.Atoi(conf.SysConfig.Port)
			if err != nil || tcpPort > 65534 || tcpPort < 1025 {
				return fmt.Errorf("%s: TCP Port '%s' is not a number between 1025 and 65535",
					StrERROR, conf.SysConfig.Port)
			}
		}
	}
	configuration = conf
	return nil
}

// GetConfiguration: returns the Configuration object
func GetConfiguration() model.Configuration {
	return configuration
}
