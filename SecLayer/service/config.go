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
	RedisQueueName   string
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
	TokenPasswd      string

	LogPath  string
	LogLevel string

	WriteGoroutineNum int
	ReadGoroutinNum   int

	SecProductInfoMap map[int]*SecProductInfoConf

	HandleUserGoroutineNum  int
	MaxRequestTimeout       int
	Read2HandleChanSize     int
	SendToWriteChanTimeout  int
	Handle2WriteChanSize    int
	SendToHandleChanTimeout int
}

type SecLayerContext struct {
	proxy2LayerRedisPool *redis.Pool
	layer2ProxyRedisPool *redis.Pool
	etcdClient           *clientv3.Client
	RwSecProductLock     sync.RWMutex
	secLayerConf         *SecLayerConf

	waitGroup        sync.WaitGroup
	Read2HandlerChan chan *SecRequest
	Handle2WriteChan chan *SecResponse

	HistoryMap     map[int]*UserBuyHistory
	HistoryMapLock sync.Mutex

	productCountMgr *ProductCountMgr
}

type SecProductInfoConf struct {
	ProductId    int
	StartTime    int64
	EndTime      int64
	Status       int
	Total        int
	Left         int
	BuyRate      float64
	SoldMaxLimit int
	secLimit     *SecLimit
	//单人限制购买的数量
	OnePersonBuyLimit int
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

type SecResponse struct {
	ProductId int
	UserId    int
	Token     string
	TokenTime int64
	Code      int
}
