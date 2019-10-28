package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

var (
	secSkillConf *SecSkillConf
)

func InitService(serviceConf *SecSkillConf) {
	secSkillConf = serviceConf
	logs.Debug("init service conf success, config:%v", secSkillConf)
}

func SecInfo(productId int) (data map[string]interface{}, code int, err error) {
	secSkillConf.RWSecProductLock.Lock()
	defer secSkillConf.RWSecProductLock.Unlock()

	v, ok := secSkillConf.SecProductInfoMap[productId]

	if !ok {
		code = ErrInvalidRequest
		err = fmt.Errorf("not found product_id:%d", productId)
		return
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTime
	data["status"] = v.Status

	return
}
