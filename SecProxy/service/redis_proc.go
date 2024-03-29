package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func WriteHandle() {
	for {
		req := <- secSkillConf.SecReqChan
		conn := secKillConf.proxy2LayerRedisPool.Get()

		data, err := json.Marshal(req)
		if err != nil {
			logs.Error("json.Marshal failed,error:%v, req:%v", err, req)
			conn.Close()
			continue
		}

		_, err = conn.Do("LPUSh", "sec_queue", string(data))
		if err != nil {
			logs.Error("lpush failed, error:%v, req:%v", err, req)
			conn.Close()
			continue
		}
		conn.Close()
	}
}

func ReadHandle() {
	for {
		conn := secKillConf.proxy2LayerRedisPool.Get()

		reply, err := conn.Do("RPOP", "recv_queue")
		if err != nil {
			logs.Error("rpop failed, error:%v", err)
			conn.Close()
			continue
		}
		data, err := redis.String(reply, err)
		if err == redis.ErrNil {
			time.Sleep(time.Second)
			conn.Close()
			continue
		}
		logs.Debug("rpop from redis succ, data:%s", string(data))


		var result SecResult
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			logs.Error("json.Unmarshal failed, error:%v", err)
			conn.Close()
			continue
		}

		userKey := fmt.Sprintf("%s_%s", result.UserId, result.ProductId)

		secKillConf.UserConnMapLock.Lock()
		resultChan, ok := secKillConf.UserConnMap[userKey]
		secKillConf.UserConnMapLock.Unlock()
		if !ok {
			conn.Close()
			logs.Warn("user not found:%v", userKey)
			continue
		}
		resultChan <- &result
		conn.Close()
	}
}
