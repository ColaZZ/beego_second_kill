package service

import (
	"fmt"
	"sync"
)

type SecLimitMgr struct {
	UserLimitMap map[int]*SecLimit
	lock sync.Mutex
}

var (
	secLimitMgr *SecLimitMgr
)

func antiSpam(req *SecRequest) (err error) {
	secLimitMgr.lock.Lock()
	limit, ok := secLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		limit := & SecLimit{
			count:   0,
			curTime: 0,
		}
		secLimitMgr.UserLimitMap[req.UserId] = limit
	}
	secIdCount := limit.Count(req.AccessTime.Unix())
	if secIdCount > secSkillConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}
	secLimitMgr.lock.Unlock()
	return 
}