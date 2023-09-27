package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"

	"testing"
	"time"
)

func TestMustNewClient(t *testing.T) {
	client, _ := MustNewClient(context.Background(), Config{
		Context:        nil,
		Endpoints:      []string{"127.0.0.1:2379"},
		EnableResolver: false,
	})

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	MustRegistry(ctx, client, Service{
		Project:      "oss",
		Name:         "file-rpc",
		MsgType:      "registry",
		Scheme:       "http",
		Host:         "127.0.0.1",
		Port:         9001,
		Replica:      "1",
		RegTime:      time.Now(),
		AuthExcludes: nil,
	})

	key := "resolver/oss/file-rpc"

	discoverService(client, key)

}

func discoverService(client *clientv3.Client, serviceName string) {
	fmt.Printf("target service %s\n", serviceName)
	key := fmt.Sprintf("%s", serviceName)

	resp, err := client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("get:%s:%s\n", ev.Key, ev.Value)
	}
}
