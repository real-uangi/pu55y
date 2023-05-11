package test

import (
	"github.com/gin-gonic/gin"
	"github.com/real-uangi/pu55y/api"
	"github.com/real-uangi/pu55y/runner"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	server := runner.Prepare()
	server.AddApi(api.GET, "/statistic/visit", func(context *gin.Context) {
		context.JSON(200, "hello")
	})
	runner.Run()
	time.Sleep(time.Hour)
}
