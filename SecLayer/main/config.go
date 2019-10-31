package main

import (
	"SecondKill/SecLayer/service"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

var (
	appConfig *service.SecLayerConf
)

func initConfig(confType, filename string) (err error) {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Println("new config failed, err:%v", err)
		return
	}

	appConfig = &service.SecLayerConf{}

	//读取日志库配置
	appConfig.LogLevel = conf.String("logs::log_level")
	if len(appConfig.LogLevel) == 0 {
		appConfig.LogLevel = "Debug"
	}

	appConfig.LogPath = conf.String("logs::log_path")
	if len(appConfig.LogPath) == 0 {
		appConfig.LogPath = "./logs"
	}

	//读取Proxy2Layer redis 相关配置
	appConfig.Proxy2LayerRedis.RedisAddr = conf.String("redis::redis_proxy2layer_addr")
	if len(appConfig.Proxy2LayerRedis.RedisAddr) == 0 {
		logs.Error("read redis_proxy2layer_addr failed")
		err = fmt.Errorf("read redis redis_proxy2layer_addr failed")
		return
	}
	appConfig.Proxy2LayerRedis.RedisMaxIdle, err = conf.Int("redis::redis_proxy2layer_idle")
	if err != nil {
		logs.Error("read redis_proxy2layer_idle failed")
		return
	}
	appConfig.Proxy2LayerRedis.RedisIdleTimeout, err = conf.Int("redis::redis_proxy2layer_idle_timeout")
	if err != nil {
		logs.Error("read redis_proxy2layer_idle_timeout failed")
		return
	}
	appConfig.Proxy2LayerRedis.RedisMaxActive, err = conf.Int("redis::redis_proxy2layer_active")
	if err != nil {
		logs.Error("read redis_proxy2layer_active failed")
		return
	}

	//读取各类goroutine线程数量
	appConfig.WriteGoroutineNum, err = conf.Int("service::write_proxy2layer_goroutine_num")
	if err != nil {
		logs.Error("read write_proxy2layer_goroutine_num failed")
		return
	}

	appConfig.ReadGoroutinNum, err = conf.Int("service::read_layer2proxy_goroutine_num")
	if err != nil {
		logs.Error("read read_layer2proxy_goroutine_num failed")
		return
	}

	return
}
