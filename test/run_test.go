package test

import (
	"github.com/gin-gonic/gin"
	"github.com/real-uangi/pu55y/api"
	"github.com/real-uangi/pu55y/plog"
	"github.com/real-uangi/pu55y/runner"
	"strconv"
	"sync"
	"testing"
)

func TestRun(t *testing.T) {
	r := runner.Prepare()
	r.AddApi(api.GET, "/statistic/visit", func(context *gin.Context) {
		context.JSON(200, "hello")
	})
	r.Run()
}

var mu sync.Mutex

func TestMu(t *testing.T) {
	runner.Prepare()
	loopLock(0)
}

func loopLock(i int) {
	mu.Lock()
	plog.Info(strconv.Itoa(i))
	if i > 3 {
		return
	}
	mu.Unlock()
	loopLock(i + 1)
}
