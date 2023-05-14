// Package runner
// @author uangi
// @date 2023/5/10 17:06
package runner

import (
	"github.com/gin-gonic/gin"
	"github.com/real-uangi/pu55y/api"
	"github.com/real-uangi/pu55y/config"
	"github.com/real-uangi/pu55y/datasource"
	"github.com/real-uangi/pu55y/plog"
	"github.com/real-uangi/pu55y/rdb"
)

type Runner struct {
	conf   *config.Configuration
	server *api.Server
}

func Prepare() *Runner {
	config.Reload()
	conf := config.GetConfig()
	runner := Runner{}
	s := api.Server{}
	s.ListenPort(conf.Http.Port)
	runner.server = &s
	runner.conf = conf
	return &runner
}

func (runner *Runner) Run() {
	if runner.conf.Datasource != nil {
		datasource.InitDataSource(&runner.conf.Datasource)
	} else {
		plog.Warn("Server running without datasource")
	}
	if &runner.conf.Redis != nil {
		rdb.Init(&runner.conf.Redis)
	} else {
		plog.Warn("Server running without redis")
	}
	plog.Info("Runner init completed")
	err := runner.server.Run()
	if err != nil {
		plog.Error(err.Error())
	}
}

func (runner *Runner) AddApi(method api.Method, uri string, processor gin.HandlerFunc) {
	runner.server.AddApi(method, uri, processor)
}
