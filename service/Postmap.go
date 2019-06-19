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
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sockmapd/base"
	"sockmapd/model"
)

func GetPostmap(request model.Request) string {
	postmapDbMap := base.GetPostmapDb()
	if len(postmapDbMap[request.Service].Key) == 0 {
		log.Printf("Service '%s' not found in configuration file", request.Service)
		return base.SmapTEMP + base.StrInternalERROR
	}
	err := postmapDbMap[request.Service].DbConn.Ping()
	if err != nil {
		log.Printf("%s: %s", base.StrERROR, err)
		// If the DB connection fails, then try to re-initialize the DB connections, this is
		// necessary if we have DB replication or cluster (not supported by the sql-driver) so
		// that connection is attempted to another available node
		err = base.InitializeDatabase()
		postmapDbMap = base.GetPostmapDb()
		if err != nil {
			log.Printf("%s: %s - after DB re-initialization", base.StrERROR, err)
			return base.SmapTEMP + base.StrTemporaryERROR
		}
	}
	var queryValue string
	if len(postmapDbMap[request.Service].Value) > 0 {
		queryValue = postmapDbMap[request.Service].Value
		if len(postmapDbMap[request.Service].Reason) > 0 {
			queryValue = queryValue + "," + postmapDbMap[request.Service].Reason
		}
	} else {
		queryValue = postmapDbMap[request.Service].Key
	}
	postmapQuery := fmt.Sprintf("SELECT %s FROM %s WHERE %s = '%s'", queryValue, postmapDbMap[request.Service].Table,
		postmapDbMap[request.Service].Key, request.Key)

	row := postmapDbMap[request.Service].DbConn.QueryRow(postmapQuery)
	var rowModel model.Response
	var response string
	if len(postmapDbMap[request.Service].Reason) > 0 {
		err = row.Scan(&rowModel.Value, &rowModel.Reason)
	} else {
		err = row.Scan(&rowModel.Value)
	}
	if err == sql.ErrNoRows {
		response = base.SmapNOTFOUND
	} else {
		if postmapDbMap[request.Service].Value == "" || rowModel.Value == base.StrOK || rowModel.Value == base.StrREJECT {
			response = base.SmapOK + rowModel.Value
		} else {
			response = rowModel.Value
		}
		if len(rowModel.Reason) > 0 {
			response = response + " " + rowModel.Reason
		}
	}
	return response
}
