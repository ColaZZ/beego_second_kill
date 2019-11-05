package main

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "SecondKill/SecAdmin/router"
)

func main() {
	err := initAll()
	if err != nil {
		panic(fmt.Errorf(""))
	}
	beego.Run()
}
