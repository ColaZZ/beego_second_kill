package service

import (
	"encoding/json"
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

	secLayerContext.layer2ProxyRedisPool, err = initRedisPool(conf.Layer2ProxyRedis)
	if err != nil {
		logs.Error("init layer2proxy redis pool redis failed, err:%v", err)
		return
	}
	return
}

func RunProcess() (err error) {
	for i := 0; i < secLayerContext.secLayerConf.ReadGoroutinNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleReader()
	}

	for i := 0 ; i < secLayerContext.secLayerConf.WriteGoroutineNum; i ++ {
		secLayerContext.waitGroup.Add(1)
		//go HandleWriter()
	}

	for i := 0; i< secLayerContext.secLayerConf.HandleUserGoroutineNum; i ++ {
		secLayerContext.waitGroup.Add(1)
		//go HandlerUser()
	}

	logs.Debug("all process goroutine started")
	secLayerContext.waitGroup.Wait()
	logs.Debug("wait all goroutine exited")
	return
}

func HandleReader() {
	for {
		conn := secLayerContext.proxy2LayerRedisPool.Get()
		for {
			data, err := redis.String(conn.Do("blpop", "queue_name", 0))
			if err != nil {
				logs.Error("")
				break
			}

			var req SecRequest
			err = json.Unmarshal([]byte(data), &req)
			if err != nil {
				logs.Error("unmarshal failed, err :%v", err)
				continue
			}

			now := time.Now().Unix()
			if now - req.AccessTime.Unix() >= int64(secLayerContext.secLayerConf.MaxRequestTimeout) {
				logs.Warn("req[%v] is expired", req)
				continue
			}

			secLayerContext.Read2HandlerChan <- &req
		}
		conn.Close()
	}
}
