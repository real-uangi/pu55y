// Package runner
// @author uangi
// @date 2023/5/10 17:06
package runner

import (
	"github.com/real-uangi/pu55y/api"
	"github.com/real-uangi/pu55y/config"
	"github.com/real-uangi/pu55y/datasource"
	"github.com/real-uangi/pu55y/rdb"
)

var conf *config.Configuration

var server *api.Server

func Prepare() *api.Server {
	config.Reload()
	conf = config.GetConfig()
	s := api.Server{}
	s.ListenPort(conf.Http.Port)
	server = &s
	return server
}

func Run() {
	if conf.Datasource != nil {
		datasource.InitDataSource(&conf.Datasource)
	}
	if &conf.Redis != nil {
		rdb.Init(&conf.Redis)
	}

}
