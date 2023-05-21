// Package config @author uangi 2023-05
package config

type Configuration struct {
	Sys        sys          `json:"sys"`
	Http       http         `json:"http"`
	Datasource []datasource `json:"datasource"`
	Redis      redis        `json:"rdb"`
	Snowflake  snowflake    `json:"snowflake"`
}

type sys struct {
	Date date
}

type http struct {
	Port string `json:"port"`
	Log  bool   `json:"log"`
}

type datasource struct {
	Name     string `json:"name"`
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

type snowflake struct {
	Lazy   bool   `json:"lazy"`
	PreGen preGen `json:"pre-generate"`
}

type preGen struct {
	Enable   bool `json:"enable"`
	Max      int  `json:"max"`
	Min      int  `json:"min"`
	Steps    int  `json:"steps"`
	Interval int  `json:"interval"`
}

type date struct {
	Zone   string
	Layout string
}
