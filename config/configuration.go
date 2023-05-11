// Package config @author uangi 2023-05
package config

type Configuration struct {
	Http       http         `json:"http"`
	Datasource []Datasource `json:"datasource"`
	Redis      Redis        `json:"rdb"`
}

type http struct {
	Port string `json:"port"`
	Log  bool   `json:"log"`
}

type Datasource struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Redis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
	PoolMin  int    `json:"poolMin"`
	PoolMax  int    `json:"poolMax"`
}
