package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	etcd_client "go.etcd.io/etcd/clientv3"
	"time"
)

var (
	redisPool  *redis.Pool
	etcdClient *etcd_client.Client
)

func initSec() (err error) {
	err = initRedis()
	if err != nil {
		logs.Error("init redis failed, err:%v", err)
		return
	}

	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed, err :%v", err)
		return
	}

	err = initLogger()
	if err != nil {
		logs.Error("init log failed, err :%v", err)
		return
	}

	err = loadSecConf()
	if err != nil {
		logs.Error("load sec conf failed, err :%v", err)
		return
	}
	return
}

func loadSecConf() (err error) {
	// key := fmt.Sprintf("%s/product", secKillConf.EtcdConf.EtcdSecProductKey)
	resp, err := etcdClient.Get(context.Background(), secKillConf.EtcdConf.EtcdSecProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err :%v", secKillConf.EtcdConf.EtcdSecProductKey, err)
		return
	}

	for k, v := range resp.Kvs {
		logs.Debug("key[%v] valued[%v]", k, v)
	}

	return
}

func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace

	}
	return logs.LevelDebug
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = secKillConf.LogPath
	config["level"] = convertLogLevel(secKillConf.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err: ", err)
		return
	}
	_ = logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}

func initRedis() (err error) {
	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisConf.RedisAddr)
		},
		MaxIdle:     secKillConf.RedisConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisConf.RedisIdleTimeout) * time.Second,
	}

	conn := redisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err :%v", err)
		return
	}
	return
}

func initEtcd() (err error) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{secKillConf.EtcdConf.EtcdAddr},
		DialTimeout: time.Duration(secKillConf.EtcdConf.Timeout) * time.Second,
	})

	if err != nil {
		logs.Error("connect etcd failed, err :%v", err)
		return
	}

	etcdClient = cli
	return
}
