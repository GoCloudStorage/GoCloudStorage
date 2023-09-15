package logic

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestGrpcServer(t *testing.T) {
	dial, _ := grpc.Dial("localhost:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := storage.NewStorageClient(dial)
	resp, err := client.GenerateDownloadURL(context.Background(), &storage.GenerateDownloadURLReq{Hash: "123456", Expire: int64(time.Second * 150)})
	fmt.Println(resp, err)
}
