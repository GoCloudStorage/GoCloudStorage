package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/opt"

	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/random"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage_engine"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/model"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type StorageServer struct {
	engine storage_engine.IStorage
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
	return fmt.Sprintf("http://%s:%s/storage/download/%s", opt.Cfg.StorageService.Host, opt.Cfg.StorageService.Port, key)
}

func getCodeKey(code string) string {
	return fmt.Sprintf("storage:code:%s", code)
}

func getFileKey(fileID string) string {
	return fmt.Sprintf("storage:file:%s", fileID)

}
