package service

import (
	"github.com/garyburd/redigo/redis"
	"go.etcd.io/etcd/clientv3"
	"sync"
)

var (
	secLayerContext = &SecLayerContext{}
)

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisIdleTimeout int
	RedisMaxActive   int
}

type EtcdConf struct {
	EtcdAddr          string
	Timeout           int
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type SecLayerConf struct {
	Proxy2LayerRedis RedisConf
	Layer2ProxyRedis RedisConf
	EtcdConfig       EtcdConf

	LogPath  string
	LogLevel string

	WriteGoroutineNum int
	ReadGoroutinNum   int

	SecProductInfoMap map[int]*SecProductInfoConf
}

type SecLayerContext struct {
	proxy2LayerRedisPool *redis.Pool
	layer2ProxyRedisPoll *redis.Pool
	etcdClient           *clientv3.Client
	RwSecProductLock     sync.Mutex
	secLayerConf         *SecLayerConf
}

type SecProductInfoConf struct {
	ProductId int
	StartTime int64
	EndTime   int64
	Status    int
	Total     int
	Left      int
}
