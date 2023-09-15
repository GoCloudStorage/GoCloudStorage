package xrpc

import (
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
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

func InitStorageRPC() (*MidClient[storage.StorageClient], error) {
	return GetGrpcClient(
		Config{
			Domain:          opt.Cfg.StorageRPC.Domain,
			Endpoints:       opt.Cfg.StorageRPC.Endpoints,
			BackoffInterval: 0,
			MaxAttempts:     0,
		},
		storage.NewStorageClient,

		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
}

func sss() {

}

func InitUserRPCClient[T any](fn func(grpc.ClientConnInterface) T) (*MidClient[T], error) {
	client, err := GetGrpcClient(Config{
		Domain:          opt.Cfg.StorageRPC.Domain,
		Endpoints:       opt.Cfg.UserRPC.Endpoints,
		BackoffInterval: 0,
		MaxAttempts:     0,
	},
		fn,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	return client, err
}
