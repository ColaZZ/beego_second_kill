package activity

import (
	"SecondKill/SecAdmin/model"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type ActivityController struct {
	beego.Controller
}

func (p *ActivityController) CreateActivity() {
	p.TplName = "activity/create.html"
	p.Layout = "layout/layout.html"
	return
}

func (p *ActivityController) ListActivity() {
	activityModel := model.NewActivityModel()
	activityList, err := activityModel.GetActivityList()
	if err != nil {
		logs.Warn("get activity list failed, err:%V", err)
		return
	}

	p.Data["activity_list"] = activityList
	p.TplName = "activity/list.html"
	p.Layout = "layout/layout.html"
	return
}
