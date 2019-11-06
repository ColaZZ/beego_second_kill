package model

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type Activity struct {
	ActivityId   int    `db:"id"`
	ActivityName string `db:"name"`
	ProductId    int    `db:"product_id"`
	StartTime    int64  `db:"start_time"`
	EndTime      int64  `db:"end_time"`
	Total        int    `db:"total"`
	Status       int    `db:"status"`

	StartTimeStr string
	EndTimeStr   string
	StatusStr    string
}

type ActivityModel struct {
}

func NewActivityModel() *ActivityModel {
	return &ActivityModel{}
}

func (p *ActivityModel) GetActivityList() (activityList []*Activity, err error) {
	sqlStr := "select id,name,product_id,start_time,end_time,total,status from activity order by id desc"
	err = Db.Select(&activityList, sqlStr)
	if err != nil {
		logs.Error("select activity from database failed, err:%v", err)
		return
	}
	for _, v := range activityList {
		t := time.Unix(v.StartTime, 0)
		v.StartTimeStr = t.Format("2006-01-02 15:04:05")

		t = time.Unix(v.EndTime, 0)
		v.EndTimeStr = t.Format("2006-01-02 15:04:05")

		now := time.Now().Unix()
		if now > v.EndTime {
			v.StatusStr = "已结束"
			continue
		}

		if v.Status == ActivityStatusNormal {
			v.StatusStr = "正常"
		} else if v.Status == ActivityStatusDisable {
			v.StatusStr = "已禁用"
		}
	}
	logs.Debug("get activity success, activityList is [%v]", activityList)
	return
}

func (p *ActivityModel) CreateActivity(activity *Activity) (err error) {
	valid, err := p.ProductValid(activity.ProductId, activity.Total)
	if err != nil {
		logs.Error("product exists falied, err:%v", err)
		return
	}
	if !valid {
		err = fmt.Errorf("product_id [%v], err:%v", activity.ProductId, err)
		logs.Error(err)
		return
	}

	if activity.StartTime <= 0 || activity.EndTime <= 0 {
		err = fmt.Errorf("invalid start_time[%v]|end_time[%v]", activity.StartTime, activity.EndTime)
		logs.Error(err)
		return
	}

	if activity.EndTime <= activity.StartTime {
		err = fmt.Errorf("start_time[%v] is greater than end_time[%v]", activity.StartTime, activity.EndTime)
		logs.Error(err)
		return
	}

	now := time.Now().Unix()
	if activity.EndTime <= now || activity.StartTime <= now {
		err = fmt.Errorf("start_time[%v]|end_time[%v] is less than now[%v]", activity.StartTime, activity.EndTime,
			now)
		logs.Error(err)
		return
	}

	sqlStr := "insert into activity(name,product_id,start_time,end_time,total values(?,?,?,?,?)"
	_, err = Db.Exec(sqlStr, activity.ActivityName, activity.ProductId, activity.StartTime, activity.EndTime,
		activity.Total)
	if err != nil {
		logs.Warn("insert into activity failed, err:%v", err)
		return
	}
	logs.Debug("insert into database succ")

	return
}

func (p *ActivityModel) ProductValid(productId, total int) (valid bool, err error) {
	sqlStr := "select id,name,total,status from product where id=?"
	var productList []*Product
	err = Db.Select(&productList, sqlStr, productId)
	if err != nil {
		logs.Warn("select product failed, err:%v", err)
		return
	}

	if len(productList) == 0 {
		err = fmt.Errorf("product[%v] is not existed", productId)
		return
	}
	if total > productList[0].Total {
		err = fmt.Errorf("product[%v]的数量非法", productId)
		return
	}
	valid = true
	return
}
