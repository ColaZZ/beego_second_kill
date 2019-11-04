package service

import (
	"github.com/garyburd/redigo/redis"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
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

	HandleUserGoroutineNum int
	MaxRequestTimeout      int
	Read2HandleChanSize    int
}

type SecLayerContext struct {
	proxy2LayerRedisPool *redis.Pool
	layer2ProxyRedisPool *redis.Pool
	etcdClient           *clientv3.Client
	RwSecProductLock     sync.Mutex
	secLayerConf         *SecLayerConf
	waitGroup            sync.WaitGroup
	Read2HandlerChan     chan *SecRequest
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
	Nance         string
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientRefence string
	//CloseNotify   <-chan bool

	//ResultChan chan *SecResult
}
