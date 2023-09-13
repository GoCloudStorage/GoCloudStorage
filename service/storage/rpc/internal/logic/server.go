package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/random"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage_engine"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/model"
	redis2 "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type StorageServer struct {
	Engine storage_engine.IStorage
	storage.UnimplementedStorageServer
}

func (s *StorageServer) CreateStorage(ctx context.Context, in *storage.CreateStorageReq) (*storage.CreateStorageResp, error) {
	parseToken, err := token.ParseUploadToken(in.Token)
	if err != nil {
		return nil, err
	}
	si := new(model.StorageInfo)
	//验参
	err = si.FindStorageByHash(parseToken.Hash)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	//创建新存储
	if errors.Is(err, gorm.ErrRecordNotFound) {
		si.Hash = parseToken.Hash
		err = si.CreateStorage()
		if err != nil {
			return nil, err
		}

		return &storage.CreateStorageResp{
			StorageId: int64(si.StorageId),
		}, nil
	}
	return nil, errors.New(response.RPC_PARAM_ERROR)
}

func (s *StorageServer) FindStorageByHash(ctx context.Context, req *storage.FindStorageByHashReq) (*storage.FindStorageByHashResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageServer) GenerateDownloadURL(ctx context.Context, req *storage.GenerateDownloadURLReq) (*storage.GenerateDownloadURLResp, error) {
	var key string
	key = random.GenerateRandomString(16)
	for redis.Client.Get(context.Background(), "storage:download:url"+key).Err() != redis2.Nil {
		key = random.GenerateRandomString(16)
	}
	s.Engine.GenerateObjectURL("storage:download:url"+key, req.Hash, time.Duration(req.Expire))
	return &storage.GenerateDownloadURLResp{URL: fmt.Sprintf("162.14.115.114:8000/storage/download/%s", key)}, nil
}
