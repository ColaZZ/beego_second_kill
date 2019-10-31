package service

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func loafProductFromEtcd(conf *SecLayerConf) (err error) {
	logs.Debug("start getting from etcd success")
	ctx, cancel:= context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	resp, err := secLayerContext.etcdClient.Get(ctx, conf.EtcdConfig.EtcdSecProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err:%v", conf.EtcdConfig.EtcdSecProductKey, err)
		return
	}
	logs.Debug("get from etcd succ, resp:%v", resp)

	var secProductInfo []SecProductInfoConf
	for k, v := range resp.Kvs {
		logs.Debug("key [%s] values [%s]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("unmarshal sec product info failed, err:%v", err)
			return
		}
		logs.Debug("sec info conf is %v", secProductInfo)
	}

	updateSecProductInfo(conf, secProductInfo)
	logs.Debug("update product info success, data:%v", secProductInfo)

	initSecProductWatcher(conf)
	logs.Debug("init etcd watcher success")
	return
}

func updateSecProductInfo(conf *SecLayerConf, secProductInfo []SecProductInfoConf) {
	var tmp map[int]*SecProductInfoConf = make(map[int]*SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		productinfo := v
		//productinfo.secLimit = &SecLimit{}
		tmp[v.ProductId] = &productinfo
	}
	secLayerContext.RwSecProductLock.Lock()
	conf.SecProductInfoMap = tmp
	secLayerContext.RwSecProductLock.Unlock()
}

func initSecProductWatcher(conf *SecLayerConf) {
	go watchSecProductKey(conf)
}

func watchSecProductKey(conf *SecLayerConf) {
	key := conf.EtcdConfig.EtcdSecProductKey
	logs.Error("begin watch key:%s", key)
	var err error
	for {
		rch := secLayerContext.etcdClient.Watch(context.Background(), key)
		var secProductInfo []SecProductInfoConf
		var getConnSuccess = true

		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s]'s config is deleted", key)
					continue
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						logs.Error("key[%s] unmarshal, err:%s", key, err)
						getConnSuccess = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q: %q ", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConnSuccess {
				logs.Debug("get config from etcd success, %v", secProductInfo)
				updateSecProductInfo(conf, secProductInfo)
			}
		}
	}
}
