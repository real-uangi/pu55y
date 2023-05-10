// Package redis @author uangi 2023-05
package redis

import (
	"context"
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
func getClient() *redis.Client {
	if client == nil {
		if clientLock.TryLock() {
			client = redis.NewClient(option)
			defer clientLock.Unlock()
		} else {
			for {
				time.Sleep(10 * time.Millisecond)
				if client != nil {
					return client
				}
			}
		}
	}
	return client
}

// Init 启动 建议放在main函数内
func Init(addr string, psw string, db int) {
	option = &redis.Options{
		Addr:         addr, //连接地址
		Password:     psw,  //密码
		DB:           db,   //库
		PoolSize:     8,    //连接池大小8
		MinIdleConns: 2,    //最小空闲 2
	}
	getClient()
}

// Set 设置
func Set(key string, value interface{}) {
	SetWithExpire(key, value, DefaultTtl)
}

// SetWithExpire 设置带过期
func SetWithExpire(key string, value interface{}, dur time.Duration) {
	err := getClient().Set(ctx, key, value, dur).Err()
	if err != nil {
		plog.Error(err.Error())
	}
}

// Get 获取
func Get(key string) string {
	res, err := getClient().Get(ctx, key).Result()
	if err != nil {
		plog.Error(err.Error())
	}
	return res
}

// SetExpire 设置过期时间
func SetExpire(key string, ttl time.Duration) bool {
	err := getClient().Expire(ctx, key, ttl).Err()
	if err != nil {
		plog.Error(err.Error())
		return false
	}
	return true
}

// TryLock 分布式锁
func TryLock(key string, parse string, ttl int) bool {
	script := redis.NewScript(`
		if redis.call('setnx', KEYS[1], ARGV[1]) == 1 
		then redis.call('pexpire', KEYS[1], tonumber(ARGV[2]));
		return 1 
		else return 0 end;
	`)
	keys := []string{key}
	args := []interface{}{parse, ttl}
	res, err := script.Run(ctx, getClient(), keys, args).Int()
	if err != nil {
		plog.Error(err.Error())
	}
	return res == 1
}

// Unlock 解锁
func Unlock(key string, parse string) {
	script := redis.NewScript(`
		if redis.call('get', KEYS[1]) == ARGV[1] 
		then return redis.call('del', KEYS[1])
		else return 0 end;
	`)
	keys := []string{key}
	args := []interface{}{parse}
	err := script.Run(ctx, getClient(), keys, args).Err()
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
