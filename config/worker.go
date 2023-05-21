// Package config @author uangi 2023-05
package config

import (
	"encoding/json"
	"os"
)

var conf = &Configuration{}

func Reload() {
	f, err := os.Open("./conf.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dc := json.NewDecoder(f)
	err = dc.Decode(conf)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Configuration {
	return conf
}
