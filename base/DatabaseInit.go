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
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sockmapd/model"
)

var postmapDb map[string]model.PostmapDb

func InitializeDatabase() error {
	conf := GetConfiguration()
	dbconnsMap := make(map[string]*sql.DB)
	postmapDbMap := make(map[string]model.PostmapDb)
	var myPostmapDb model.PostmapDb
	for i := range conf.Postmaps {
		if dbconnsMap[conf.Postmaps[i].Database] == nil {
			db, err := openDatabaseConnection(conf.Database, conf.Postmaps[i].Database)
			if err != nil {
				log.Printf("%s: %s", StrERROR, err)
				return fmt.Errorf("%s: %s for database %s", StrERROR, err, conf.Postmaps[i].Database)
			}
			dbconnsMap[conf.Postmaps[i].Database] = db
		}
		myPostmapDb.DbConn = dbconnsMap[conf.Postmaps[i].Database]
		myPostmapDb.Table  = conf.Postmaps[i].Table
		myPostmapDb.Key    = conf.Postmaps[i].Key
		myPostmapDb.Value  = conf.Postmaps[i].Value
		myPostmapDb.Reason = conf.Postmaps[i].Reason
		postmapDbMap[conf.Postmaps[i].Service] = myPostmapDb
	}
	postmapDb = postmapDbMap
	return nil
}

func openDatabaseConnection(dbConfig model.DbConfig, database string) (*sql.DB, error){
	// MySQL
	var db *sql.DB
	// Try all hosts in the Host array and save the connection that first succeeds
	for i := range dbConfig.Host {
		dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbConfig.Username,
			dbConfig.Password, dbConfig.Host[i], dbConfig.Port, database)

		db, err := sql.Open("mysql", dbUri)
		if err != nil {
			log.Printf("%s: %s", StrERROR, err)
			continue
		}
		err = db.Ping()
		if err != nil {
			log.Printf("%s: %s", StrERROR, err)
			continue
		}
		return db, nil
	}
	return db, error(fmt.Errorf("could not get a DB connection"))
}

func GetPostmapDb() map[string]model.PostmapDb {
	return postmapDb
}