package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"sync"
)

type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	lock         sync.Mutex
}

func antiSpam(req *SecRequest) (err error) {
	_, ok := secSkillConf.idBlackMap[req.UserId]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Error("userid[%v] is blocked by id black list", req.UserId)
		return
	}

	_, ok = secSkillConf.ipBlackMap[req.ClientAddr]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Error("userId[%v] ip[%v] is blocked by ip black list", req.UserId, req.ClientAddr)
		return
	}

	// user id 频率控制
	secSkillConf.secLimitMgr.lock.Lock()
	limit, ok := secSkillConf.secLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		limit := &Limit{
			secLimit: &SecLimit{},
			minLimit: &MinLimit{},
		}
		secKillConf.secLimitMgr.UserLimitMap[req.UserId] = limit
	}

	secIdCount := limit.secLimit.Count(req.AccessTime.Unix())
	minIdCount := limit.minLimit.Count(req.AccessTime.Unix())

	// ip 频率控制
	limit, ok = secKillConf.secLimitMgr.IpLimitMap[req.ClientAddr]
	if !ok {
		limit := &Limit{
			secLimit: &SecLimit{},
			minLimit: &MinLimit{},
		}
		secSkillConf.secLimitMgr.IpLimitMap[req.ClientAddr] = limit
	}
	secIpCount := limit.secLimit.Count(req.AccessTime.Unix())
	minIpCount := limit.minLimit.Count(req.AccessTime.Unix())
	secSkillConf.secLimitMgr.lock.Unlock()


	if secIpCount > secKillConf.AccessLimitConf.IPSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	if minIpCount > secKillConf.AccessLimitConf.IPMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	if secIdCount > secKillConf.AccessLimitConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	if minIdCount > secKillConf.AccessLimitConf.UserMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}
	return
}
