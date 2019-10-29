package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"

	"crypto/md5"
)

var (
	secSkillConf *SecSkillConf
)

func InitService(serviceConf *SecSkillConf) {
	secSkillConf = serviceConf
	logs.Debug("init service conf success, config:%v", secSkillConf)
}

func SecInfo(productId int) (data []map[string]interface{}, code int, err error) {
	secSkillConf.RWSecProductLock.Lock()
	defer secSkillConf.RWSecProductLock.Unlock()

	item, code, err := SecInfoById(productId)
	if err != nil {
		return
	}
	data = append(data, item)
	return
}

func SecInfoList() (data []map[string]interface{}, code int, err error) {
	secSkillConf.RWSecProductLock.RLock()
	defer secSkillConf.RWSecProductLock.RUnlock()

	for _, v := range secSkillConf.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			logs.Error("get product id failed, product id :%d, err: %v", v.ProductId, err)
			continue
		}
		logs.Debug("get product [%d], result [%v], all[%v] v:[%v]",
			v.ProductId, item, secSkillConf.SecProductInfoMap, v)
		data = append(data, item)
	}
	return
}

func SecInfoById(productId int) (data map[string]interface{}, code int, err error) {
	secSkillConf.RWSecProductLock.RLock()
	defer secSkillConf.RWSecProductLock.Lock()

	v, ok := secSkillConf.SecProductInfoMap[productId]

	if !ok {
		code = ErrInvalidRequest
		err = fmt.Errorf("not found product_id:%d", productId)
		return
	}

	start := false
	end := false
	status := "success"

	now := time.Now().Unix()
	//秒杀未开始
	if now-v.StartTime < 0 {
		start = false
		end = false
		status = "sec kill do not start"
		code = ErrActiveNotStart
	}

	//秒杀开始
	if now-v.StartTime > 0 {
		start = true
	}

	// 秒杀结束
	if now-v.EndTime > 0 {
		start = false
		end = true
		status = "sec kill has ended"
		code = ErrActiveAlreadyEnd
	}

	// 售罄
	if v.Status == ProductStatusForceSaleOut || v.Status == ProductStatusSaleOut {
		start = false
		end = true
		status = "product sale out"
		code = ErrActiveSaleOut
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start"] = start
	data["end"] = end
	data["status"] = status

	return
}

func SecKill(req *SecRequest) (data []map[string]interface{}, code int, err error) {
	secSkillConf.RWSecProductLock.RLock()
	defer secSkillConf.RWSecProductLock.RUnlock()

	err = userCheck(req)
	if err != nil {
		code = ErrUserCheckAuthFailed
		logs.Warn("userId[%s] invalid, check failed, req[%v]", req.UserId, req)
		return
	}

	err = antiSpam(req)
	if err != nil {
		code = ErrUserServiceBusy
		logs.Warn("userId[%s] invalid, check failed, req[%v]", req.UserId, req)
		return
	}
	return
}

func userCheck(req *SecRequest) (err error) {
	authData := fmt.Sprintf("%d:%s", req.UserId, secSkillConf.CookieSecretKey)
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))
	if authSign != req.UserAuthSign {
		err = fmt.Errorf("invalid user cookie auth")
		return
	}
	return
}