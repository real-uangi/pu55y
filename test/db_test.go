package test

import (
	"github.com/real-uangi/pu55y/datasource"
	"testing"
)

func TestDB(t *testing.T) {
	datasource.InitDataSource("192.168.0.211", "5432", "", "", "")
}
