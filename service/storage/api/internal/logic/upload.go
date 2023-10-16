package logic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
	Key          string
	ChunkNumber  int
	ChunksNumber int
}

type ContentRange struct {
	Start int
	End   int
	Total int
}

func getUploadChunkSizeKey(key string, chunk int) string {
	return fmt.Sprintf("storage:upload:%s:chunk:%d:size", key, chunk)
}

func UploadPart(f *bytes.Reader, uploadReq *UploadReq) (*model.StorageInfo, error) {
	var (
		object model.StorageInfo
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
		object = model.StorageInfo{
			Hash: uploadReq.Key,
			Size: 0,
		}
		if err = object.CreateStorage(); err != nil {
			logrus.Errorf("failed create object record, err: %v", err)
			return nil, err
		}
	}

	// 保存数据，并将块大小存储在redis中
	chunkSizeKey := getUploadChunkSizeKey(uploadReq.Key, uploadReq.ChunkNumber)
	s, err := redis.Get(context.Background(), chunkSizeKey)
	if err != nil && err != redis.Nil {
		logrus.Error("redis get size err:", err)
		return nil, err
	}
	if s == "" {
		s = "0"
	}
	size, err := strconv.Atoi(s)
	if err != nil {
		logrus.Error("strconv.Atoi err：", err)
		return nil, err
	}
	object.Size = size
	//文件长度
	lenth := f.Len()

	//该块已经完成了，直接返回
	if size == uploadReq.ContentRange.Total {
		return &object, nil
	}

	//保存
	if err = local.Client.SaveChunk(uploadReq.Key, uploadReq.ChunkNumber, f, int64(uploadReq.ContentRange.Start)); err != nil {
		logrus.Errorf("failed save local, err: %v", err)
		return nil, err
	}
	//上传完成
	redis.SetEx(context.Background(), chunkSizeKey, size+lenth, time.Minute*30)
	object.Size = size + lenth

	//该块全部上传，记录完成的快号
	if object.Size == uploadReq.ContentRange.Total {
		redis.Client.SAdd(context.Background(), getUploadFinishChunkKey(uploadReq.Key), uploadReq.ChunkNumber)
	}

	return &object, err
}

func getUploadFinishChunkKey(key string) string {
	return fmt.Sprintf("storage:upload:%s:chunck", key)
}
