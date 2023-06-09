package test

import (
	"github.com/real-uangi/pu55y/plog"
	"github.com/real-uangi/pu55y/runner"
	"github.com/real-uangi/pu55y/snowflake"
	"strconv"
	"sync"
	"testing"
)

var idMap = make(map[string]int)

var succeed int

var al sync.Mutex

func TestSnowflake(t *testing.T) {
	succeed = 0
	runner.Prepare()
	a := 100
	for i := 0; i < a; i++ {
		go gen(i)
	}
	for succeed < a {
	}
}

func gen(i int) {
	plog.Warn("group " + strconv.Itoa(i) + " started")
	for j := 0; j < 100; j++ {
		plog.Info(snowflake.NextId().String())
	}
	succeed++
	plog.Info("group " + strconv.Itoa(i) + " finished!")
}
