package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/local"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/mq"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/api/internal/logic"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

func getUploadChunkKey(key string) string {
	return fmt.Sprintf("storage:upload:%s:lockChunk", key)
}

func getUploadFinishChunkKey(key string) string {
	return fmt.Sprintf("storage:upload:%s:chunck", key)
}

func (a *API) upload(c *fiber.Ctx) error {
	var (
		uploadToken string
	)
	type resp struct {
		StorageId uint64 `json:"storage_id"`
	}

	uploadToken = c.Params("token")
	parseUploadToken, err := token.ParseUploadToken(uploadToken)
	if err != nil {
		logrus.Error("parse upload token err:", err)
		return response.Resp400(c, nil)
	}

	// 解析请求头
	uploadReq, err := parasUploadHeader(c)
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}
	if uploadReq.Key != parseUploadToken.Key {
		return response.Resp400(c, nil)
	}

	//分布式锁
	flag := redis.SetLock(context.Background(), getUploadChunkKey(uploadReq.Key), time.Minute*30)
	if !flag {
		return response.Resp202(c, nil)
	}
	defer redis.ReleaseLock(context.Background(), getUploadChunkKey(uploadReq.Key))

	//验证分块是否已被人上传
	flag = redis.Client.SIsMember(context.Background(), getUploadFinishChunkKey(uploadReq.Key), uploadReq.ChunkNumber).Val()
	if flag {
		return response.Resp200(c, nil)
	}

	//上传
	object, err := logic.UploadPart(bytes.NewReader(c.Body()), uploadReq)
	if err != nil {
		return response.Resp500(c, nil)
	}

	//上传完成
	redis.Client.SAdd(context.Background(), getUploadChunkKey(uploadReq.Key), uploadReq.ChunkNumber)
	// 上传到指定大小，合并

	if object.Size == uploadReq.TotalSize {
		path, err := local.Client.MergeChunk(object.Hash, object.Size)
		if err != nil {
			return response.Resp500(c, nil, fmt.Sprintf("merge chunk failed, err: %v", err))
		}
		object.RealPath = path

		if err = object.UpdateStorage(); err != nil {
			return response.Resp500(c, nil, fmt.Sprintf("save object record failed, err: %v", err))
		}
		//完成上传，通知异步
		if opt.Cfg.StorageRPC.IsRemote == 1 {
			marshal, err := json.Marshal(object.StorageId)
			if err != nil {
				logrus.Error("异步上传oss marshal storageId err:", err)
				return response.Resp200(c, resp{StorageId: object.StorageId}, "上传完成，marshal id err")
			}
			err = mq.Publish("", "transfer-task", marshal)
			if err != nil {
				logrus.Error("异步上传oss err:", err)
				return response.Resp200(c, resp{StorageId: object.StorageId}, "上传完成，异步上传失败")
			}
		}
		return response.Resp200(c, resp{StorageId: object.StorageId}, "上传完成，合并成功")
	}
	return response.Resp200(c, nil)

}
