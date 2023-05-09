package test

import (
	"github.com/real-uangi/pu55y/plog"
	"testing"
)

func TestLog(t *testing.T) {
	plog.Info("This is an info log")
	plog.Warn("This is a warning log")
	plog.Error("This is an error log")
}
