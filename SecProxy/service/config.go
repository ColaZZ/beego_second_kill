package service

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

const (
	ProductStatusNormal       = 0
	ProductStatusSaleOut      = 1
	ProductStatusForceSaleOut = 2
)

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr          string
	Timeout           int
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type SecSkillConf struct {
	RedisBlackConf       RedisConf
	EtcdConf           EtcdConf
	LogPath            string
	LogLevel           string
	SecProductInfoMap  map[int]*SecProductInfoConf
	RWSecProductLock   sync.RWMutex
	CookieSecretKey    string
	UserSecAccessLimit int
	ReferWhiteList     []string
	IPSecAccessLimit   int

	ipBlackMap map[string]bool
	idBlackMap map[int]bool

	blackRedisPool *redis.Pool
}

type SecProductInfoConf struct {
	ProductId int
	StartTime int64
	EndTime   int64
	Status    int
	Total     int
	Left      int
}

type SecRequest struct {
	ProductId     int
	Source        string
	AuthCode      string
	SecTime       string
	Sign          string
	Nance         string
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientRefence string
}
