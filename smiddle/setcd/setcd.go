package setcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"sync"
	"time"
)

type EtcdClient struct {
	EndPoints string
	Client    *clientv3.Client
	KV        clientv3.KV
	Lock      sync.Mutex
}

// NewEtcdClient 新建一个Etcd实例
func NewEtcdClient(etcdEndPoints string, username, password string) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdEndPoints, ","),
		DialTimeout: 5 * time.Second,
		Username:    username,
		Password:    password,
	})

	if err != nil {
		return nil, err
	}

	etcdClient := new(EtcdClient)
	etcdClient.Client = client
	etcdClient.KV = clientv3.NewKV(client)
	etcdClient.EndPoints = etcdEndPoints
	return etcdClient, nil
}
