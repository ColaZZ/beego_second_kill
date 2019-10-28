package main

import (
	"SecondKill/SecProxy/service"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strings"
)

var (
	secKillConf = &service.SecSkillConf{
		SecProductInfoMap: make(map[int]*service.SecProductInfoConf, 1024),
	}
)

func initConfig() (err error) {
	redisAddr := beego.AppConfig.String("redis_addr")
	etcdAddr := beego.AppConfig.String("etcd_addr")

	logs.Error("read redis config success :%v", redisAddr)
	logs.Error("read etcd config successs :%v", etcdAddr)

	secKillConf.RedisConf.RedisAddr = redisAddr
	secKillConf.EtcdConf.EtcdAddr = etcdAddr

	if len(redisAddr) == 0 || len(etcdAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s] or etcd[%s] config is null", redisAddr, etcdAddr)
		return
	}

	redis_max_idle, err := beego.AppConfig.Int("redis_max_idle")
	if err != nil {
		err = fmt.Errorf("init config fialed, read redis_max_idle failed, err :%v", err)
		return
	}

	redis_max_active, err := beego.AppConfig.Int("redis_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_max_active failed, err :%v", err)
		return
	}

	redis_idle_timeout, err := beego.AppConfig.Int("redis_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_idle_timeout failed, err :%v", err)
		return
	}

	secKillConf.RedisConf.RedisMaxIdle = redis_max_idle
	secKillConf.RedisConf.RedisMaxActive = redis_max_active
	secKillConf.RedisConf.RedisIdleTimeout = redis_idle_timeout

	logPath := beego.AppConfig.String("log_path")
	logLevel := beego.AppConfig.String("log_level")

	logs.Error("read log path config success :%v", logPath)
	logs.Error("read log level config success :%v", logLevel)

	secKillConf.LogLevel = logLevel
	secKillConf.LogPath = logPath

	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_timeout error:%v", err)
		return
	}
	secKillConf.EtcdConf.Timeout = etcdTimeout
	secKillConf.EtcdConf.EtcdSecKeyPrefix = beego.AppConfig.String("etcd_sec_key_prefix")
	if len(secKillConf.EtcdConf.EtcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("init config failed, read etcd_sec_key error:%v", err)
		return
	}

	productKey := beego.AppConfig.String("etcd_product_key")
	if len(productKey) == 0 {
		err = fmt.Errorf("init config failed, read etcd_product_key error:%v", err)
		return
	}

	if strings.HasSuffix(secKillConf.EtcdConf.EtcdSecKeyPrefix, "/") == false {
		secKillConf.EtcdConf.EtcdSecKeyPrefix = secKillConf.EtcdConf.EtcdSecKeyPrefix + "/"
	}

	secKillConf.EtcdConf.EtcdSecProductKey = fmt.Sprintf("%s%s", secKillConf.EtcdConf.EtcdSecKeyPrefix, productKey)
	

	return

}
