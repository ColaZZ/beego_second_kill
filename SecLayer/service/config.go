package service

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisIdleTimeout int
	RedisMaxActive     int
}

type EtcdConf struct {
	EtcdAddr          string
	Timeout           int
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type SecLayerConf struct {
	Proxy2LayerRedis  RedisConf
	Layer2ProxyRedis  RedisConf
	EtcdConfig        EtcdConf

	LogPath           string
	LogLevel          string

	WriteGoroutineNum int
	ReadGoroutinNum   int
}
