package test

import (
	"github.com/real-uangi/pu55y/plog"
	"github.com/real-uangi/pu55y/redis"
	"strconv"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	redis.Init("192.168.0.211:6379", "", 0)
	redis.Set("hello", "hi")
	res := redis.Get("hello")
	plog.Info(res)
	plog.Info(strconv.FormatInt(redis.Incr("incr", 1), 10))
	time.Sleep(time.Second)
}
