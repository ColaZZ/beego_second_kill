package main

import (
	"SecondKill/SecProxy/service"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	secKillConf = &service.SecSkillConf{
		RedisConf: service.RedisConf{},
		EtcdConf:  service.EtcdConf{},
		LogPath:   "",
		LogLevel:  "",
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

	return

}
