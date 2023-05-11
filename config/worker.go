// Package config @author uangi 2023-05
package config

import (
	"encoding/json"
	"github.com/real-uangi/pu55y/plog"
	"os"
)

var conf = &Configuration{}

func Reload() {
	f, err := os.Open("../conf.json")
	plog.TryThrow(err)
	defer f.Close()
	dc := json.NewDecoder(f)
	err = dc.Decode(conf)
	plog.TryThrow(err)
}

func GetConfig() *Configuration {
	return conf
}
