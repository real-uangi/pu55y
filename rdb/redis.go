// Package rdb @author uangi 2023-05
package rdb

import (
	"context"
	"github.com/real-uangi/pu55y/config"
	"github.com/real-uangi/pu55y/plog"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var option *redis.Options

var clientLock sync.Mutex

var client *redis.Client

var ctx = context.Background()

const (
	DefaultTtl = 24 * time.Hour
)

// GetClient 获取已有的连接
func GetClient() *redis.Client {
	if client == nil {
		if clientLock.TryLock() {
			defer clientLock.Unlock()
			if client == nil {
				client = redis.NewClient(option)
			}
		} else {
			for {
				time.Sleep(5 * time.Millisecond)
				return GetClient()
			}
		}
	}
	return client
}

// Init 启动 建议放在main函数内
func Init(c *config.Redis) {
	option = &redis.Options{
		Addr:         c.Addr,     //连接地址
		Password:     c.Password, //密码
		DB:           c.Db,       //库
		PoolSize:     c.PoolMax,  //连接池大小8
		MinIdleConns: c.PoolMin,  //最小空闲 2
	}
	GetClient()
	if client != nil {
		plog.Info("Redis client connected to " + c.Addr)
		ping := client.Ping(ctx)
		if ping != nil {
			r, e := ping.Result()
			if e != nil {
				plog.Error(e.Error())
				plog.Error("Redis connection check failed")
			} else {
				plog.Info("Redis 'ping' successfully responded as : " + r)
			}
		}
	}
}

// Set 设置
func Set(key string, value interface{}) {
	SetWithExpire(key, value, DefaultTtl)
}

// SetWithExpire 设置带过期
func SetWithExpire(key string, value interface{}, dur time.Duration) {
	err := GetClient().Set(ctx, key, value, dur).Err()
	if err != nil {
		plog.Error(err.Error())
	}
}

// Get 获取
func Get(key string) string {
	res, err := GetClient().Get(ctx, key).Result()
	if err != nil {
		plog.Error(err.Error())
	}
	return res
}

// SetExpire 设置过期时间
func SetExpire(key string, ttl time.Duration) bool {
	err := GetClient().Expire(ctx, key, ttl).Err()
	if err != nil {
		plog.Error(err.Error())
		return false
	}
	return true
}

// TryLock 分布式锁
func TryLock(key string, parse string, ttl int) bool {
	script := redis.NewScript(`
		if rdb.call('setnx', KEYS[1], ARGV[1]) == 1 
		then rdb.call('pexpire', KEYS[1], tonumber(ARGV[2]));
		return 1 
		else return 0 end;
	`)
	keys := []string{key}
	args := []interface{}{parse, ttl}
	res, err := script.Run(ctx, GetClient(), keys, args).Int()
	if err != nil {
		plog.Error(err.Error())
	}
	return res == 1
}

// Unlock 解锁
func Unlock(key string, parse string) {
	script := redis.NewScript(`
		if rdb.call('get', KEYS[1]) == ARGV[1] 
		then return rdb.call('del', KEYS[1])
		else return 0 end;
	`)
	keys := []string{key}
	args := []interface{}{parse}
	err := script.Run(ctx, GetClient(), keys, args).Err()
	if err != nil {
		plog.Error(err.Error())
	}
}

// Incr 自增
func Incr(key string, count int64) int64 {
	res, err := client.IncrBy(ctx, key, count).Result()
	if err != nil {
		plog.Error(err.Error())
	}
	return res
}
