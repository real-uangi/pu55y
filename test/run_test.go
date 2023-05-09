package test

import (
	"github.com/real-uangi/pu55y/config"
	"github.com/real-uangi/pu55y/plog"
	"testing"
)

func TestRun(t *testing.T) {

	config.Reload()
	plog.Info(config.GetConfig().Http.Port)

}
