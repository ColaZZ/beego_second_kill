package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var (
	secKillConf *SecSkillConf
	//redisPool  *redis.Pool
)

func InitService(serviceConf *SecSkillConf) (err error) {
	secSkillConf = serviceConf
	logs.Debug("init service conf success, config:%v", secSkillConf)
	err = loadBlackList()
	if err != nil {
		logs.Error("load black err:%v", err)
		return
	}
	return
}

func initBlackRedis() (err error) {
	secKillConf.blackRedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisBlackConf.RedisAddr)
		},
		MaxIdle:     secKillConf.RedisBlackConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisBlackConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisBlackConf.RedisIdleTimeout) * time.Second,
	}

	conn := secKillConf.blackRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err :%v", err)
		return
	}
	return
}

func loadBlackList() (err error) {
	err = initBlackRedis()
	if err != nil {
		logs.Error("init black reids failed, err:%v", err)
		return
	}

	conn := secSkillConf.blackRedisPool.Get()
	defer conn.Close()
	// id black list
	reply, err := conn.Do("hgetall", "idblacklist")
	idList, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed, err:%v", err)
		return
	}
	for _, v := range idList {
		id, err := strconv.Atoi(v)
		if err != nil {
			logs.Warn("invalid user, id[%v]", id)
			continue
		}
		secSkillConf.idBlackMap[id] = true
	}

	// ip list
	reply, err = conn.Do("hgetall", "idblacklist")
	ipList, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed, err:%v", err)
		return
	}
	for _, v := range ipList {
		secSkillConf.ipBlackMap[v] = true
	}
	return

}
