package model

import (
	"github.com/jmoiron/sqlx"
	"go.etcd.io/etcd/clientv3"
)

var (
	Db             *sqlx.DB
	EtcdClient     *clientv3.Client
	EtcdPrefix     string
	EtcdProductKey string
)

func Init(db *sqlx.DB, etcdClient *clientv3.Client, prefix, productKey string) (err error) {
	Db = db
	EtcdClient = etcdClient
	EtcdPrefix = prefix
	EtcdProductKey = productKey
	return
}
