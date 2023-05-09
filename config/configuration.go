package config

type Configuration struct {
	Http       http         `json:"http"`
	Datasource []datasource `json:"datasource"`
	Redis      redis        `json:"redis"`
}

type http struct {
	Port string `json:"port"`
	Log  bool   `json:"log"`
}

type datasource struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type redis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
	PoolMin  int    `json:"poolMin"`
	PoolMax  int    `json:"poolMax"`
}
