package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var (
	secKillConf *SecSkillConf
)

func InitService(serviceConf *SecSkillConf) (err error) {
	secSkillConf = serviceConf
	logs.Debug("init service conf success, config:%v", secSkillConf)
	err = loadBlackList()
	if err != nil {
		logs.Error("load black err:%v", err)
		return
	}
	logs.Debug("init service succ, config:%v", secKillConf)

	err = initProxy2LayerRedis()
	if err != nil {
		logs.Error("load proxy2layer redis failed, err:%v", err)
		return
	}

	secKillConf.secLimitMgr = &SecLimitMgr{
		UserLimitMap: make(map[int]*Limit, 10000),
		IpLimitMap:   make(map[string]*Limit, 10000),
	}

	secKillConf.SecReqChan = make(chan *SecRequest, secKillConf.SecReqChanSize)
	secKillConf.UserConnMap = make(map[string]chan *SecResult, 10000)

	initRedisProcessFunc()
	return
}

func initRedisProcessFunc() {
	for i := 0; i < secSkillConf.WriteProxy2LayerGoroutineNum; i++ {
		go WriteHandle()
	}

	for i := 0; i < secSkillConf.ReadProxy2LayerGoroutineNum; i ++ {
		go ReadHandle()
	}
}

func initProxy2LayerRedis() (err error) {
	secKillConf.proxy2LayerRedisPool = &redis.Pool{
		MaxIdle:     secKillConf.RedisProxy2LayerConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisProxy2LayerConf.RedisAddr)
		},
	}

	conn := secKillConf.proxy2LayerRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}

func initLayer2ProxyRedis() (err error) {
	secKillConf.layer2ProxyRedisPool = &redis.Pool{
		MaxIdle:     secKillConf.RedisLayer2ProxyConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisLayer2ProxyConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisLayer2ProxyConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisLayer2ProxyConf.RedisAddr)
		},
	}

	conn := secKillConf.layer2ProxyRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
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

	go SyncIpBlackList()
	go SyncIdBlackList()
	return
}

func SyncIpBlackList() {
	var ipList []string
	lastTime := time.Now().Unix()

	for {
		conn := secSkillConf.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackiplist", time.Second)
		ip, err := redis.String(reply, err)
		if err != nil {
			continue
		}

		curTime := time.Now().Unix()
		ipList = append(ipList, ip)

		if len(ipList) == 100 || curTime-lastTime > 5 {
			secSkillConf.RWBlackLock.Lock()
			for _, v := range ipList {
				secSkillConf.ipBlackMap[v] = true
			}
			secSkillConf.RWBlackLock.Unlock()
			lastTime = curTime
			logs.Info("sync ip list from redis success, ip [%v]", ipList)
		}
	}
}

func SyncIdBlackList() {
	for {
		conn := secSkillConf.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackidlist", time.Second)
		id, err := redis.Int(reply, err)
		if err != nil {
			continue
		}

		secSkillConf.RWBlackLock.Lock()
		secSkillConf.idBlackMap[id] = true
		secSkillConf.RWBlackLock.Unlock()

		logs.Info("sync id list from redis succ, ip[%v]", id)
	}
}
