package activity

import (
	"SecondKill/SecAdmin/model"
	"fmt"
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

func (p *ActivityController) SubmitActivity() {
	activityModel := model.NewActivityModel()
	var activity model.Activity

	p.TplName = "activity/list.html"
	p.Layout = "layout/layout.html"

	var err error
	var Error string = "success"
	defer func() {
		if err != nil {
			p.Data["Error"] = Error
			p.TplName = "activity/error.html"
		}
	}()

	name := p.GetString("activity_name")
	if len(name) == 0 {
		Error = "活动的名字不能为空"
		err = fmt.Errorf("activity name can not be null")
		return
	}

	productId, err := p.GetInt("product_id")
	if err != nil {
		err = fmt.Errorf("商品id非法, err:%v", err)
		Error = err.Error()
		return
	}

	startTime, err := p.GetInt64("start_time")
	if err != nil {
		err = fmt.Errorf("开始时间非法, err:%v", err)
		Error = err.Error()
		return
	}

	endTime, err := p.GetInt64("end_time")
	if err != nil {
		err = fmt.Errorf("结束时间非法, err:%v", err)
		Error = err.Error()
		return
	}

	total, err :=p.GetInt("Total")
	if err != nil {
		err = fmt.Errorf("商品数量非法, err:%v", err)
		Error = err.Error()
		return
	}

	activity.ActivityName = name
	activity.ProductId = productId
	activity.StartTime = startTime
	activity.EndTime = endTime
	activity.Total = total

	err = activityModel

}
