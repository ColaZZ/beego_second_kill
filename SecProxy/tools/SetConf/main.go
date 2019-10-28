package SetConf

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type SecInfoConf struct {
	ProductId int
	StartTime int
	EndTime   int
	Status    int
	Total     int
	Left      int
}

var (
	EtcdKey = "/sk/backend/seckill/product"
)

func SetLogConfToEtcd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:22379", "localhost:2379", "localhost:33279"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect etcd failed,err:", err)
		return
	}

	fmt.Println("coonect etcd success")
	defer func() {
		_ = cli.Close()
	}()

	var SecInfoConfArry []SecInfoConf
	SecInfoConfArry = append(
		SecInfoConfArry,
		SecInfoConf{
			ProductId: 1029,
			StartTime: 1505008800,
			EndTime:   1505012400,
			Status:    0,
			Total:     1000,
			Left:      1000,
		},
	)
	SecInfoConfArry = append(
		SecInfoConfArry,
		SecInfoConf{
			ProductId: 1027,
			StartTime: 1505008800,
			EndTime:   1505012400,
			Status:    0,
			Total:     2000,
			Left:      1000,
		},
	)

	data, err := json.Marshal(SecInfoConfArry)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, EtcdKey, string(data))
	if err != nil {
		fmt.Println("etcd put %s failed", EtcdKey)
		return
	}
	cancel()

	context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	if err != nil {
		fmt.Println("etcd get failed", EtcdKey)
		return
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

func main() {
	SetLogConfToEtcd()
}
