package logic

import (
	"context"
	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/snowflake"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/model"
	"gorm.io/gorm"
)

type storageServer struct {
	storage.UnimplementedStorageServer
}

func (s *storageServer) CreateStorage(ctx context.Context, in *storage.CreateStorageReq) (*storage.CreateStorageResp, error) {
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
		//雪花算法生成ID
		id, err := snowflake.GetID()
		if err != nil {
			return nil, err
		}

		si.StorageId = int64(id)
		si.Hash = parseToken.Hash
		err = si.CreateStorage()
		if err != nil {
			return nil, err
		}

		return &storage.CreateStorageResp{
			StorageId: si.StorageId,
		}, nil
	}
	return nil, errors.New(response.PARAM_ERROR)
}
