package xrpc

import (
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitRPCClient[T any](fn func(grpc.ClientConnInterface) T) (*MidClient[T], error) {
	client, err := GetGrpcClient(Config{
		Domain:          opt.Cfg.StorageRPC.Domain,
		Endpoints:       opt.Cfg.StorageRPC.Endpoints,
		BackoffInterval: 0,
		MaxAttempts:     0,
	},
		fn,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	return client, err
}
