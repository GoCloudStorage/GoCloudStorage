package rpc

import (
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var UserClient user.UserServiceClient

func ClientInit() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := user.NewUserServiceClient(conn)

	UserClient = client
}
