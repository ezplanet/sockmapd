package model

import "database/sql"

type PostmapDb struct {
	DbConn		*sql.DB
	Table		string
	Key			string
	Value 		string
	Reason		string
}
