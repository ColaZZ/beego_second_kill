package service

import (
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func InitSecLayer(conf *SecLayerConf) (err error) {
	err = initRedis(conf)
	if err != nil {
		logs.Error("init redis failed, err:%v", err)
		return
	}
	logs.Debug("init redis success")

	err = initEtcd(conf)
	if err != nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	}
	logs.Debug("init etcd success")

	err = loafProductFromEtcd(conf)
	if err != nil {
		logs.Error("load product from etcd failed, err:%v", err)
		return
	}

	secLayerContext.secLayerConf = conf

	return
}

func initEtcd(conf *SecLayerConf) (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            []string{conf.EtcdConfig.EtcdAddr},
		DialTimeout:          time.Duration(conf.EtcdConfig.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:%v", err)
		return
	}

	secLayerContext.etcdClient = cli
	logs.Debug("init etcd success")
	return
}