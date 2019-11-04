package main

import (
	"SecondKill/SecLayer/service"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"strings"
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
	appConfig.Proxy2LayerRedis.RedisQueueName = conf.String("redis::redis_proxy2layer_queue_name")
	if len(appConfig.Proxy2LayerRedis.RedisQueueName) == 0 {
		logs.Error("read redis_proxy2layer_queue_name failed")
		err = fmt.Errorf("read redis redis_proxy2layer_queue_name failed")
		return
	}

	//读取Layer2Proxy redis 相关配置
	appConfig.Layer2ProxyRedis.RedisAddr = conf.String("redis::redis_layer2proxy_addr")
	if len(appConfig.Layer2ProxyRedis.RedisAddr) == 0 {
		logs.Error("read redis_layer2proxy_addr failed")
		err = fmt.Errorf("read redis redis_layer2proxy_addr failed")
		return
	}
	appConfig.Layer2ProxyRedis.RedisMaxIdle, err = conf.Int("redis::redis_layer2proxy_idle")
	if err != nil {
		logs.Error("read redis_layer2proxy_idle failed")
		return
	}
	appConfig.Layer2ProxyRedis.RedisIdleTimeout, err = conf.Int("redis::redis_layer2proxy_idle_timeout")
	if err != nil {
		logs.Error("read redis_layer2proxy_idle_timeout failed")
		return
	}
	appConfig.Layer2ProxyRedis.RedisMaxActive, err = conf.Int("redis::redis_layer2proxy_active")
	if err != nil {
		logs.Error("read redis_layer2proxy_active failed")
		return
	}
	appConfig.Layer2ProxyRedis.RedisQueueName = conf.String("redis::redis_layer2proxy_queue_name")
	if len(appConfig.Layer2ProxyRedis.RedisQueueName) == 0 {
		logs.Error("read redis_layer2proxy_queue_name failed")
		err = fmt.Errorf("read redis redis_layer2proxy_queue_name failed")
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

	appConfig.HandleUserGoroutineNum, err = conf.Int("service::handle_user_goroutine_num")
	if err != nil {
		logs.Error("read handle_user_goroutine_num failed")
		return
	}

	appConfig.MaxRequestTimeout, err = conf.Int("service::max_request_wait_timeout")
	if err != nil {
		logs.Error("read max_request_wait_timeout failed")
		return
	}

	appConfig.Read2HandleChanSize, err = conf.Int("service::read2handle_chan_size")
	if err != nil {
		logs.Error("read read2handle_chan_size failed")
		return
	}

	appConfig.Handle2WriteChanSize, err = conf.Int("service::handle2write_chan_size")
	if err != nil {
		logs.Error("read handle2write_chan_size failed")
		return
	}

	appConfig.SendToWriteChanTimeout, err = conf.Int("service::send_to_write_chan_timeout")
	if err != nil {
		logs.Error("read send_to_write_chan_timeout failed")
		return
	}

	appConfig.SendToHandleChanTimeout, err = conf.Int("service::send_to_handle_chan_timeout")
	if err != nil {
		logs.Error("read send_to_handle_chan_timeout failed")
		return
	}

	appConfig.TokenPasswd = conf.String("service::seckill_token_passwd")
	if len(appConfig.TokenPasswd) == 0 {
		logs.Error("read seckill_token_passwd failed")
		err = fmt.Errorf("read redis seckill_token_passwd failed")
		return
	}

	//读取etcd相关的配置

	appConfig.EtcdConfig.EtcdAddr = conf.String("etcd::server_addr")
	if len(appConfig.TokenPasswd) == 0 {
		logs.Error("read service::seckill_token_passwd failed")
		err = fmt.Errorf("read service::seckill_token_passwd failed")
		return
	}

	etcdTimeout, err := conf.Int("etcd::etcd_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_timeout error:%v", err)
		return
	}

	appConfig.EtcdConfig.Timeout = etcdTimeout
	appConfig.EtcdConfig.EtcdSecKeyPrefix = conf.String("etcd::etcd_sec_key_prefix")
	if len(appConfig.EtcdConfig.EtcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("init config failed, read etcd_sec_key error:%v", err)
		return
	}

	productKey := conf.String("etcd::etcd_product_key")
	if len(productKey) == 0 {
		err = fmt.Errorf("init config failed, read etcd_product_key error:%v", err)
		return
	}

	if strings.HasSuffix(appConfig.EtcdConfig.EtcdSecKeyPrefix, "/") == false {
		appConfig.EtcdConfig.EtcdSecKeyPrefix = appConfig.EtcdConfig.EtcdSecKeyPrefix + "/"
	}

	appConfig.EtcdConfig.EtcdSecProductKey = fmt.Sprintf("%s%s", appConfig.EtcdConfig.EtcdSecKeyPrefix, productKey)
	return
}
