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
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type StorageServer struct {
	storage.UnimplementedStorageServer
}

func (s *StorageServer) UpdateStorage(ctx context.Context, in *storage.UpdateStorageReq) (*storage.UpdateStorageResp, error) {
	if in.StorageId == 0 {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}

	si := new(model.StorageInfo)
	si.StorageId = uint64(in.StorageId)
	si.Hash = in.Hash
	if in.IsComplete != false {
		si.IsComplete = true
	}
	si.RealPath = in.RealPath
	si.Size = int(in.Size)
	err := si.UpdateStorage()
	if err != nil {
		logrus.Error("update storage err:", err)
		return nil, errors.New(response.RPC_DB_ERROR)
	}
	return &storage.UpdateStorageResp{}, nil

}

func (s *StorageServer) GetStorageByStorageId(ctx context.Context, in *storage.GetStorageByStorageIdReq) (*storage.GetStorageByStorageIdResp, error) {
	si := new(model.StorageInfo)
	si.StorageId = uint64(in.StorageId)
	err := si.GetStorageByStorageId()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error("GetStorageByStorageId err:", err)
		return nil, errors.New(response.RPC_DB_ERROR)
	}
	//未找到记录
	if err != nil {
		return &storage.GetStorageByStorageIdResp{}, nil
	}
	return &storage.GetStorageByStorageIdResp{
		StorageId:  int64(si.StorageId),
		Hash:       si.Hash,
		Size:       int32(si.Size),
		IsComplete: si.IsComplete,
		RealPath:   si.RealPath,
	}, nil
}

func (s *StorageServer) CreateStorage(ctx context.Context, in *storage.CreateStorageReq) (*storage.CreateStorageResp, error) {

	parseToken, err := token.ParseUploadToken(in.Token)
	if err != nil {
		logrus.Error("parse upload token err:", err)
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}

	si := new(model.StorageInfo)
	//验参
	err = si.FindStorageByHash(parseToken.Hash)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error("FindStorageByHash err:", err)
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}

	//创建新存储
	if errors.Is(err, gorm.ErrRecordNotFound) {
		si.Hash = parseToken.Hash
		err = si.CreateStorage()
		if err != nil {
			logrus.Error("CreateStorage err:", err)
			return nil, errors.New(response.RPC_DB_ERROR)
		}

		return &storage.CreateStorageResp{
			StorageId: int64(si.StorageId),
		}, nil
	}
	return nil, errors.New(response.RPC_PARAM_ERROR)
}

func (s *StorageServer) FindStorageByHash(ctx context.Context, in *storage.FindStorageByHashReq) (*storage.FindStorageByHashResp, error) {
	if in == nil {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}
	si := new(model.StorageInfo)
	err := si.FindStorageByHash(in.Hash)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error("find storage by hash err:", err)
		return nil, errors.New(response.RPC_DB_ERROR)
	}

	//未找到
	if err != nil {
		return &storage.FindStorageByHashResp{}, nil
	}

	return &storage.FindStorageByHashResp{
		StorageId:  int64(si.StorageId),
		Size:       int32(si.Size),
		IsComplete: si.IsComplete,
		RealPath:   si.RealPath,
	}, nil
}

func (s *StorageServer) GenerateDownloadURL(ctx context.Context, req *storage.GenerateDownloadURLReq) (*storage.GenerateDownloadURLResp, error) {
	var (
		fileKey = getFileKey(strconv.Itoa(int(req.GetFileID())))
	)
	key, err := redis.Get(ctx, fileKey)
	if errors.Is(err, redis.Nil) {
		// 没有该文件url
		code := random.GenerateRandomString(64)
		realPath, err := storage_engine.Client.GenerateObjectURL(req.GetHash(), time.Duration(req.GetExpire()))
		if err != nil {
			log.Println("222")
			return nil, err
		}
		// 生成code:fileRealPath映射
		err = redis.SetEx(ctx, getCodeKey(code), realPath, time.Duration(req.GetExpire()))
		if err != nil {
			return nil, err
		}

		// 生成file:code映射
		err = redis.SetEx(ctx, fileKey, code, time.Duration(req.GetExpire()))
		if err != nil {
			return nil, err
		}
		return &storage.GenerateDownloadURLResp{URL: generateURL(code)}, nil
	} else if err == nil {
		// 有该文件url
		return &storage.GenerateDownloadURLResp{URL: generateURL(key)}, nil
	} else {
		// 其他错误
		return nil, err
	}
}

func (s *StorageServer) GetRealPathByCode(ctx context.Context, req *storage.GetRealPathByCodeReq) (*storage.GetRealPathByCodeResp, error) {
	realPath, err := redis.Get(ctx, getCodeKey(req.GetCode()))
	if err != nil {
		return nil, err
	}
	return &storage.GetRealPathByCodeResp{Path: realPath}, nil
}

func generateURL(key string) string {
	return fmt.Sprintf("localhost:8000/storage/download/%s", key)
}

func getCodeKey(code string) string {
	return fmt.Sprintf("storage:code:%s", code)
}

func getFileKey(fileID string) string {
	return fmt.Sprintf("storage:file:%s", fileID)
}
