package main

import (
	"SecondKill/SecAdmin/model"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var Db *sqlx.DB
var EtcdClient *clientv3.Client

func initDb() (err error){
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", AppConfig.mysqlConf.Username, AppConfig.mysqlConf.Passwd,
		AppConfig.mysqlConf.Host, AppConfig.mysqlConf.Port, AppConfig.mysqlConf.Database)
	database, err := sqlx.Open("mysql", dns)
	if err != nil {
		logs.Error("open mysql failed, err:%v", err)
		return
	}
	Db = database
	logs.Debug("connect to mysql success")
	return
}

func initAll() (err error){
	err = initConfig()
	if err != nil {
		logs.Warn("init config faiied, err:%v", err)
		return
	}

	err = initDb()
	if err != nil {
		logs.Warn("init DB failed, err:%v", err)
		return
	}
	err = model.Init(Db, EtcdClient, AppConfig.etcdConf.EtcdKeyPrefix, AppConfig.etcdConf.ProductKey)
	if err != nil {
		logs.Warn("init model failed, err:%v", err)
		return
	}
	err = initEtcd()
	if err!= nil {
		logs.Warn("init etcd failed, err:%v", err)
		return
	}
	return
}

func initEtcd() (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            []string{AppConfig.etcdConf.Addr},
		DialTimeout:          time.Duration(AppConfig.etcdConf.Timeout) * time.Second,

	})
	if err != nil {
		logs.Error("connect etcd failed,err:%v", err)
		return
	}

	EtcdClient = cli
	logs.Debug("init etcd success")
	return
}
