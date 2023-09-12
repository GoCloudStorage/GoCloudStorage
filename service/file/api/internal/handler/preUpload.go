package handler

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/xrpc"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func preUpload(ctx *fiber.Ctx) error {
	type preUploadReq struct {
		UploaderId int    `json:"uploader,omitempty" form:"uploader"`
		FileName   string `json:"file_name,omitempty" form:"file_name" `
		Ext        string `json:"ext,omitempty" form:"ext"`
		Path       string `json:"path,omitempty" form:"path"`
		Hash       string `json:"hash,omitempty" form:"hash"`
		Size       int    `json:"size,omitempty" form:"size"`
	}

	type uploadResp struct {
		Token     string `json:"token,omitempty"`
		StorageId int64  `json:"storageId,omitempty"`
	}

	p := new(preUploadReq)

	if err := ctx.BodyParser(p); err != nil {
		return err
	}

	//验参
	num := p.Size/opt.Cfg.File.BlockSize + 1
	token, err := token.GenerateToken(p.Hash, num, p.Size)
	if err != nil {
		logrus.Error("GenerateToken err:", err)
		return err
	}

	//查询是否存在该存储
	client, err := xrpc.InitRPCClient(storage.NewStorageClient)
	if err != nil {
		logrus.Error("InitRPCClient err:", err)
		return err
	}

	resp, err := client.NewSession().FindStorageByHash(context.Background(), &storage.FindStorageByHashReq{Hash: p.Hash})
	if err != nil {
		logrus.Error("FindStorageByHash err:", err)
		return err
	}

	//未存在该存储，新建存储
	if resp == nil {
		createStorageResp, err := client.NewSession().CreateStorage(context.Background(), &storage.CreateStorageReq{Token: token})
		if err != nil {
			logrus.Error("CreateStorage err:", err)
			return err
		}

		createStorageResp

	} else {
		return response.Resp200(ctx, uploadResp{
			Token:     token,
			StorageId: resp.StorageId,
		})
	}
	return nil

}

func PreDownload(ctx *fiber.Ctx) error {
	panic("not impl")
}

func GetAll(ctx *fiber.Ctx) error {
	panic("not impl")
}
