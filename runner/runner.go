// Package runner
// @author uangi
// @date 2023/5/10 17:06
package runner

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/real-uangi/pu55y/api"
	"github.com/real-uangi/pu55y/config"
	"github.com/real-uangi/pu55y/datasource"
	"github.com/real-uangi/pu55y/plog"
	"github.com/real-uangi/pu55y/rdb"
	"github.com/real-uangi/pu55y/snowflake"
)

const banner = `
  _____    _    _   _____   _____  __     __
 |  __ \  | |  | | | ____| | ____| \ \   / /
 | |__) | | |  | | | |__   | |__    \ \_/ /      version: v1
 |  ___/  | |  | | |___ \  |___ \    \   /       for more information 
 | |      | |__| |  ___) |  ___) |    | |        see https://github.com/real-uangi/pu55y
 |_|       \____/  |____/  |____/     |_|   

`

type Runner struct {
	server *api.Server
}

func Prepare() *Runner {
	fmt.Print(banner)
	// configurations
	config.Reload()
	// gin server
	runner := Runner{}
	s := api.Server{}
	s.ListenPort(config.GetConfig().Http.Port)
	runner.server = &s
	// dependencies
	if config.GetConfig().Datasource != nil {
		datasource.InitDataSource(&config.GetConfig().Datasource)
	} else {
		plog.Warn("Server running without datasource")
	}
	if &config.GetConfig().Redis != nil {
		rdb.Init(&config.GetConfig().Redis)
		snowflake.Init()
	} else {
		plog.Warn("Server running without redis ,therefore snowflake id won't work")
	}
	return &runner
}

func (runner *Runner) Run() {
	plog.Info("Runner init completed")
	err := runner.server.Run()
	if err != nil {
		plog.Error(err.Error())
	}
}

func (runner *Runner) AddApi(method api.Method, uri string, processor gin.HandlerFunc) {
	runner.server.AddApi(method, uri, processor)
}
