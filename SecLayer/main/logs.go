package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["level"] = appConfig.LogLevel
	config["filename"] = appConfig.LogPath

	configStr,err := json.Marshal(config)
	if err != nil {
		fmt.Println("init logger failed, marshal:%v", err)
		return
	}
	err = logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}
