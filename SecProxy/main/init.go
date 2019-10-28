package main

import (
	"SecondKill/SecProxy/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	etcd_client "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
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

	service.InitService(secKillConf)
	initSecProductWatch()
	logs.Info("init sec succ")

	return
}

func initSecProductWatch() {
	go watchSecProductKey(secKillConf.EtcdConf.EtcdSecProductKey)
}

func watchSecProductKey(key string) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect etcd failed, err: ", err)
		return
	}

	logs.Debug("begin etcd watch key:", key)
	for {
		rch := cli.Watch(context.Background(), key)
		var getConfSucc = true
		var secProductInfo []service.SecProductInfoConf
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}
				if ev.Type == mvccpb.PUT && key == string(ev.Kv.Key) {
					err = json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						logs.Error("key %s unmarshal %s, err:%s", key, ev.Kv.Value, err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

			}

			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", secProductInfo)
				updateSecProductInfo(secProductInfo)
			}
		}
	}
}

func updateSecProductInfo(secProductInfo []service.SecProductInfoConf) {
	var tmp map[int]*service.SecProductInfoConf = make(map[int]*service.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		tmp[v.ProductId] = &v
	}
	secKillConf.RWSecProductLock.Lock()
	secKillConf.SecProductInfoMap = tmp
	secKillConf.RWSecProductLock.Unlock()
}

func loadSecConf() (err error) {
	// key := fmt.Sprintf("%s/product", secKillConf.EtcdConf.EtcdSecProductKey)
	resp, err := etcdClient.Get(context.Background(), secKillConf.EtcdConf.EtcdSecProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err :%v", secKillConf.EtcdConf.EtcdSecProductKey, err)
		return
	}

	var secProductInfo []service.SecProductInfoConf
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] valued[%v]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("Unmarshal sec product info failed, err:%v", err)
			return
		}
		logs.Debug("sec info conf is [%v]", secProductInfo)

	}

	// secKillConf.SecProductInfo

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
