package handler

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/file"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/xrpc"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func preUpload(ctx *fiber.Ctx) error {
	type preUploadReq struct {
		UploaderId int64  `json:"uploader,omitempty" form:"uploader"`
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

	num := p.Size/opt.Cfg.File.BlockSize + 1
	token, err := token.GenerateUploadToken(p.Hash, num, p.Size)
	if err != nil {
		logrus.Error("GenerateUploadToken err:", err)
		return err
	}

	//验参
	fileClient, err := xrpc.InitRPCClient(file.NewFileClient)
	if err != nil {
		logrus.Error("InitRPCClient err", err)
		return err
	}
	//查询文件是否存在
	info, err := fileClient.NewSession().FindFileByUserIdAndFileInfo(context.Background(), &file.FindFileByUserIdAndFileInfoReq{
		UserId:   p.UploaderId,
		Path:     p.Path,
		FileName: p.FileName,
		Ext:      p.Ext,
	})
	if err != nil {
		logrus.Error("FindFileByUserIdAndFileInfo err:", err)
		return err
	}

	//已存在该文件，直接返回存储id
	if info != nil {
		return response.Resp200(ctx, uploadResp{
			Token:     token,
			StorageId: info.StorageId,
		})

	}

	//该用户未存在该文件，查看存储桶是否有通用的该数据
	//查询是否存在该存储
	storageClient, err := xrpc.InitRPCClient(storage.NewStorageClient)
	if err != nil {
		logrus.Error("InitRPCClient err:", err)
		return err
	}

	findStorageResp, err := storageClient.NewSession().FindStorageByHash(context.Background(), &storage.FindStorageByHashReq{Hash: p.Hash})
	if err != nil {
		logrus.Error("FindStorageByHash err:", err)
		return err
	}

	//未存在该存储，新建存储
	if findStorageResp == nil {
		createStorageResp, err := storageClient.NewSession().CreateStorage(context.Background(), &storage.CreateStorageReq{Token: token})
		if err != nil {
			logrus.Error("CreateStorage err:", err)
			return err
		}

		//新建用户文件信息
		_, err = fileClient.NewSession().CreateFile(context.Background(), &file.CreateFileReq{
			UserId:    p.UploaderId,
			Path:      p.Path,
			FileName:  p.FileName,
			Ext:       p.Ext,
			Hash:      p.Hash,
			Size:      int32(p.Size),
			BlockSize: int32(opt.Cfg.File.BlockSize),
			StorageId: createStorageResp.StorageId,
		})

		if err != nil {
			logrus.Error("CreateFile err:", err)
			return err
		}
		return response.Resp200(ctx, uploadResp{
			Token:     token,
			StorageId: createStorageResp.StorageId,
		})

	} else { //存在该存储，直接返回存储id
		return response.Resp200(ctx, uploadResp{
			Token:     token,
			StorageId: findStorageResp.StorageId,
		})
	}

}
