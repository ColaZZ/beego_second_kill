package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type MysqlConfig struct {
	Username string
	Passwd   string
	Host     string
	Port     int
	Database string
}

type Config struct {
	mysqlConf MysqlConfig
}

var AppConfig Config

func initConfig() (err error) {
	//mysql 配置
	mysqlUsername := beego.AppConfig.String("mysql_user_name")
	if len(mysqlUsername) == 0 {
		logs.Error("read mysql_user_name faield, err", err)
		return
	}
	AppConfig.mysqlConf.Username = mysqlUsername

	mysqlPasswd := beego.AppConfig.String("mysql_passwd")
	if len(mysqlPasswd) == 0 {
		logs.Error("read mysql_passwd faield, err", err)
		return
	}
	AppConfig.mysqlConf.Passwd = mysqlPasswd

	mysqlHost := beego.AppConfig.String("mysql_host")
	if len(mysqlHost) == 0 {
		logs.Error("read mysql_host faield, err", err)
		return
	}
	AppConfig.mysqlConf.Host = mysqlHost

	mysqlDatabase := beego.AppConfig.String("mysql_database")
	if len(mysqlDatabase) == 0 {
		logs.Error("read mysql_database faield, err", err)
		return
	}
	AppConfig.mysqlConf.Database = mysqlDatabase

	mysqlPort, err := beego.AppConfig.Int("mysql_port")
	if err != nil {
		logs.Error("read mysql_port failed, err%v", err)
		return
	}
	AppConfig.mysqlConf.Port = mysqlPort



	return
}
