package handler

import (
	"bytes"
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage_engine"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage_engine/local"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const expireTime = 100 * time.Second

func (a *API) upload(ctx *fiber.Ctx) error {
	type uploadReq struct {
		StorageId int64  `json:"storage_id,omitempty"`
		Token     string `json:"token,omitempty"`
		Data      []byte `json:"data,omitempty"`
	}

	type uploadResp struct {
		IsComplete bool `json:"is_complete,omitempty"`
		Num        int  `json:"num,omitempty"` //已经完成的切片编号,从1开始
	}

	u := new(uploadReq)
	err := ctx.BodyParser(u)
	if err != nil {
		return err
	}

	uploadToken, err := token.ParseUploadToken(u.Token)
	if err != nil {
		return err
	}

	//分布式锁处理并发
	ok := redis.SetLock(context.Background(), "lock_"+uploadToken.Hash, expireTime)
	defer redis.ReleaseLock(context.Background(), "lock_"+uploadToken.Hash)

	//未获取锁，返回最新的上传进度
	if !ok {
		result, err := redis.Client.Get(context.Background(), uploadToken.Hash).Result()
		if err != nil {
			logrus.Error("get block num err:", err)
			return err
		}
		num := 0
		if result != "" {
			num, _ = strconv.Atoi(result)

		}

		return response.Resp200(ctx, &uploadResp{
			IsComplete: false,
			Num:        num,
		})
	}

	//获取了锁，开始上传
	num, err := redis.Client.Incr(context.Background(), uploadToken.Hash).Result()
	if err != nil {
		logrus.Error("redis incr block num err:", err)
		return err
	}

	var s local.StorageEngine
	err = s.UploadChunk(&storage_engine.UploadChunkRequest{
		FileMD5: uploadToken.Hash,
		Data:    bytes.NewReader(u.Data),
		PartNum: int(num),
	})

	if err != nil {
		logrus.Error("UploadChunk err:", err)
		return err
	}

	//上传成功,是否传输完所有分片
	if uploadToken.PartNum <= int(num) {
		//更新存储文件完整性
		_, err := a.storageRPCClient.UpdateStorage(context.Background(), &storage.UpdateStorageReq{
			StorageId:  u.StorageId,
			IsComplete: true,
		})
		if err != nil {
			logrus.Error("UpdateStorage err:", err)
			return err
		}

		return response.Resp200(ctx, &uploadResp{
			IsComplete: true,
			Num:        int(num),
		})
	}
	return response.Resp200(ctx, &uploadResp{
		IsComplete: false,
		Num:        int(num),
	})

}
