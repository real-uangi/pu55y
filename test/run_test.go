package test

import (
	"github.com/gin-gonic/gin"
	"github.com/real-uangi/pu55y/api"
	"github.com/real-uangi/pu55y/runner"
	"testing"
)

func TestRun(t *testing.T) {
	r := runner.Prepare()
	r.AddApi(api.GET, "/statistic/visit", func(context *gin.Context) {
		context.JSON(200, "hello")
	})
	r.Run()
}
