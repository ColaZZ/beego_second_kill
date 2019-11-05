package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

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
	return
}
