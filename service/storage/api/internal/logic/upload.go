package logic

import (
	"bytes"
	"context"
	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/local"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UploadReq struct {
	ContentRange ContentRange
	TotalSize    int

	Key         string
	ChunkNumber int
}

type ContentRange struct {
	Start int
	End   int
	Total int
	start int64
}

func UploadPart(f *bytes.Reader, uploadReq *UploadReq) (*model.StorageInfo, error) {
	var (
		object *model.StorageInfo
		err    error
		find   bool
	)
	// 获取 object record
	err = object.FirstByHash(uploadReq.Key)

	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		logrus.Errorf("failed get object, err: %v", err)
		return nil, err
	}
	if err == nil {
		find = true
	}
	//没找到,创建存储
	if !find {
		object = &model.StorageInfo{
			Hash: uploadReq.Key,
			Size: uploadReq.TotalSize,
		}
		if err = object.CreateStorage(); err != nil {
			logrus.Errorf("failed create object record, err: %v", err)
			return nil, err
		}
	}

	// 保存数据，并将大小存储在redis中
	sizeKey := getUploadSizeKey(object.Hash)
	s, err := redis.Get(context.Background(), sizeKey)
	if err != nil {
		logrus.Error("redis get size err:", err)
		return nil, err
	}
	size, err := strconv.Atoi(s)
	if err != nil {
		logrus.Error("strconv.Atoi err：", err)
		return nil, err
	}

	if err = local.Client.SaveChunk(uploadReq.Key, uploadReq.ChunkNumber, f, int64(uploadReq.ContentRange.Start)); err != nil {
		logrus.Errorf("failed save local, err: %v", err)
		return nil, err
	}
	redis.SetEx(context.Background(), sizeKey, size+f.Len(), time.Minute*30)
	object.Size = size + f.Len()

	return object, err
}

func getUploadSizeKey(key string) string {
	return "storage:upload:size:" + key
}
