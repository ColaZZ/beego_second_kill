package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func initRedisPool(redisconf RedisConf) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisconf.RedisAddr)
		},
		MaxIdle:     redisconf.RedisMaxIdle,
		MaxActive:   redisconf.RedisMaxActive,
		IdleTimeout: time.Duration(redisconf.RedisIdleTimeout) * time.Second,
	}

	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}
	return
}

func initRedis(conf *SecLayerConf) (err error) {
	secLayerContext.proxy2LayerRedisPool, err = initRedisPool(conf.Proxy2LayerRedis)
	if err != nil {
		logs.Error("init proxy2layer redis pool redis failed, err:%v", err)
		return
	}

	secLayerContext.layer2ProxyRedisPoll, err = initRedisPool(conf.Layer2ProxyRedis)
	if err != nil {
		logs.Error("init layer2proxy redis pool redis failed, err:%v", err)
		return
	}
	return
}
