package main

import (
	"context"
	"encoding/json"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/mq"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/xrpc"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/model"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	opt.InitConfig()
	pg.Init(opt.Cfg.Pg.Host, opt.Cfg.Pg.User, opt.Cfg.Pg.Password, opt.Cfg.Pg.DBName, opt.Cfg.Pg.Port)
	mq.Init(opt.Cfg.Mq.Addr, opt.Cfg.Mq.Username, opt.Cfg.Mq.Password)
	wg.Add(1)
	storageRPC, err := xrpc.GetGrpcClient(
		xrpc.Config{
			Domain:          opt.Cfg.StorageRPC.Domain,
			Endpoints:       opt.Cfg.StorageRPC.Endpoints,
			BackoffInterval: 0,
			MaxAttempts:     0,
		},
		storage.NewStorageClient,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		panic(err)
	}
	mq.Consume(&wg, "transfer-task", func(wg *sync.WaitGroup, msgs <-chan amqp091.Delivery) {
		for msg := range msgs {
			var (
				task model.Task
			)
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				logrus.Errorf("unmarshal task failed, err: %v", err)
				break
			}

			if _, err := storageRPC.NewSession().UploadOSS(context.Background(), &storage.UploadOSSReq{StorageID: task.StorageID}); err != nil {
				logrus.Error("upload object to oss failed, err: %v", err)
				continue
			}
			msg.Ack(false)
		}
		wg.Done()
	})
	wg.Wait()
}
