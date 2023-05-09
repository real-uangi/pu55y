package test

import (
	"github.com/gin-gonic/gin"
	"github.com/real-uangi/pu55y/api"
	"github.com/real-uangi/pu55y/plog"
	"testing"
	"time"
)

func TestGin(t *testing.T) {
	var server api.Server
	server.AddApi(api.GET, "/hello", hello)
	server.ListenPort(5001)
	err := server.Run()
	if err != nil {
		plog.Error(err.Error())
	}
	time.Sleep(time.Hour)
}

func hello(c *gin.Context) {
	c.JSON(200, api.Success("hello"))
}
