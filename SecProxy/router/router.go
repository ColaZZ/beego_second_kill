package router

import (
	"SecondKill/SecProxy/controller"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func init() {
	logs.Debug("enter router init")
	beego.Router("/seckill", &controller.SkillController{}, "*:SecKill")
	beego.Router("/secinof", &controller.SkillController{}, "*SecInfo")
}
