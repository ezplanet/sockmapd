package model

type Configuration struct {
	Database	DbConfig
	SysConfig	SysConfig
	Postmaps	[]Postmap
}

type DbConfig struct {
	Host		[]string
	Port		string
	Username	string
	Password	string
}

type SysConfig struct {
	Port		string
	Loglevel	string
	Logfile		string
}

type Postmap struct {
	Service		string
	Database	string
	Table		string
	Key			string
	Value 		string
	Reason		string
}