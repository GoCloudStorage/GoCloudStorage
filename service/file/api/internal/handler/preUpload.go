package handler

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (a *API) preUpload(ctx *fiber.Ctx) error {
	type preUploadReq struct {
		UploaderId int64  `json:"upload_id,omitempty" form:"upload_id"`
		FileName   string `json:"file_name,omitempty" form:"file_name" `
		Ext        string `json:"ext,omitempty" form:"ext"`
		Path       string `json:"path,omitempty" form:"path"`
		Hash       string `json:"hash,omitempty" form:"hash"`
		Size       int    `json:"size,omitempty" form:"size"`
		Expire     int64  `json:"expire,omitempty" form:"ep"`
	}
	type preUploadResp struct {
		URL      string `json:"url,omitempty"`
		ChunkNum int32  `json:"chunk_num,omitempty"`
	}

	p := new(preUploadReq)

	if err := ctx.BodyParser(p); err != nil {
		return response.Resp400(ctx, nil)
	}

	

	//将创建用户信息加入消息队列，完成上传后调用
	//将上传到云加入消息队列，完成上传后调用
	//获取下载链接
	resp, err := a.storageRPCClient.GetUploadURL(context.Background(), &storage.GetUploadURLReq{
		Hash:   p.Hash,
		Expire: p.Expire,
		Size:   0,
	})
	if err != nil {
		logrus.Error("storageRPCClient.GetUploadURL err:", err)
		return response.Resp500(ctx, nil)
	}
	return response.Resp200(ctx, preUploadResp{
		URL:      resp.Url,
		ChunkNum: resp.ChunkNum,
	})
}
