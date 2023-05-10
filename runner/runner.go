// Package runner
// @author uangi
// @date 2023/5/10 17:06
package runner

import (
	"github.com/real-uangi/pu55y/api"
	"github.com/real-uangi/pu55y/config"
	"github.com/real-uangi/pu55y/datasource"
)

var conf *config.Configuration

var server *api.Server

func Prepare() {

	config.Reload()
	conf = config.GetConfig()

	server := api.Server{}
	server.ListenPort(conf.Http.Port)

}

func Run() {

	datasource.InitDataSource(&conf.Datasource[0])

}
