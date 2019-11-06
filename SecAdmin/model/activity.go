package model

import (
	"github.com/astaxie/beego/logs"
	"time"
)

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type Activity struct {
	ActivityId   int   `db:"id"`
	ActivityName int   `db:"name"`
	ProductId    int   `db:"product_id"`
	StartTime    int64 `db:"start_time"`
	EndTime      int64 `db:"end_time"`
	Total        int   `db:"total"`
	Status       int   `db:"status"`

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
