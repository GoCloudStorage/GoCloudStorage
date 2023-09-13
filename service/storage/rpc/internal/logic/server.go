package logic

import (
	"context"
	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/model"
	"gorm.io/gorm"
)

type StorageServer struct {
	storage.UnimplementedStorageServer
}

func (s *StorageServer) FindStorageByHash(ctx context.Context, in *storage.FindStorageByHashReq) (*storage.FindStorageByHashResp, error) {
	if in == nil {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}
	si := new(model.StorageInfo)
	err := si.FindStorageByHash(in.Hash)
	if err != nil {
		return nil, errors.New(response.RPC_DB_ERROR)
	}
	return &storage.FindStorageByHashResp{
		StorageId:  int64(si.StorageId),
		Size:       int32(si.Size),
		IsComplete: si.IsComplete,
		RealPath:   si.RealPath,
	}, nil
}

func (s *StorageServer) CreateStorage(ctx context.Context, in *storage.CreateStorageReq) (*storage.CreateStorageResp, error) {
	if in == nil {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}
	parseToken, err := token.ParseUploadToken(in.Token)
	if err != nil {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}
	si := new(model.StorageInfo)
	//验参
	err = si.FindStorageByHash(parseToken.Hash)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}

	//创建新存储
	if errors.Is(err, gorm.ErrRecordNotFound) {
		si.Hash = parseToken.Hash
		err = si.CreateStorage()
		if err != nil {
			return nil, errors.New(response.RPC_DB_ERROR)
		}

		return &storage.CreateStorageResp{
			StorageId: int64(si.StorageId),
		}, nil
	}
	return nil, errors.New(response.RPC_PARAM_ERROR)
}
