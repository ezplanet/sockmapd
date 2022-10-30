package model

type Configuration struct {
	Database  DbConfig
	SysConfig SysConfig
	Postmaps  []Postmap
}

type DbConfig struct {
	Host     []string
	Port     string
	Username string
	Password string
}

type SysConfig struct {
	RunAsUser  string `json:"run_as_user"`
	RunAsGroup string `json:"run_as_group"`
	Port       string
	Loglevel   string
	Logfile    string
}

type Postmap struct {
	Service  string
	Database string
	Table    string
	Key      string
	Value    string
	Reason   string
}
