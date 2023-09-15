package handler

import (
	"bytes"
	"context"
	"errors"
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
		Data      string `json:"data,omitempty"`
		Num       int    `json:"num,omitempty"`
	}

	type uploadResp struct {
		IsComplete bool `json:"is_complete,omitempty"`
		FinishNum  int  `json:"num,omitempty"` //已经完成的切片编号,从1开始
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

	//缓存中已上传的编号
	result, err := redis.Get(context.Background(), uploadToken.Hash)
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.Error("get block num err:", err)
		return err
	}
	num := 0
	if result != "" {
		num, _ = strconv.Atoi(result)
	}

	if num != u.Num-1 {
		return response.Resp200(ctx, &uploadResp{
			IsComplete: false,
			FinishNum:  num,
		}, "已秒传该块儿")
	}

	//未获取锁，返回最新的上传进度
	if !ok {
		return response.Resp200(ctx, &uploadResp{
			IsComplete: false,
			FinishNum:  num,
		})
	}

	//获取了锁，开始上传
	var s local.StorageEngine
	err = s.UploadChunk(&storage_engine.UploadChunkRequest{
		FileMD5: uploadToken.Hash,
		Data:    bytes.NewReader([]byte(u.Data)),
		PartNum: u.Num,
	})

	if err != nil {
		logrus.Error("UploadChunk err:", err)
		return err
	}

	//上传成功
	_, err = redis.Client.Incr(context.Background(), uploadToken.Hash).Result()
	if err != nil {
		logrus.Error("redis incr block num err:", err)
		return err
	}

	//是否传输完所有分片
	if uploadToken.PartNum <= num+1 {
		//合并
		err = s.MergeChunk(uploadToken.Hash, uploadToken.PartNum, uploadToken.Size)
		if err != nil {
			logrus.Error("MergeChunk err:", err)
			return err
		}

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
			FinishNum:  u.Num,
		})
	}
	return response.Resp200(ctx, &uploadResp{
		IsComplete: false,
		FinishNum:  u.Num,
	})

}
