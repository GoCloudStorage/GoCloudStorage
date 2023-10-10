package handler

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/file"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (a *API) preUpload(ctx *fiber.Ctx) error {
	type preUploadReq struct {
		FileName  string `json:"file_name,required" form:"file_name,required" `
		Ext       string `json:"ext,required" form:"ext,required"`
		Path      string `json:"path,required" form:"path,required"`
		Hash      string `json:"hash,required" form:"hash,required"`
		Size      int32  `json:"size,required" form:"size,required"`
		IsPrivate bool   `json:"is_private" form:"is_private"`
		Expire    int64  `json:"expire" form:"expire"`
	}
	type preUploadResp struct {
		URL      string `json:"url,omitempty"`
		ChunkNum int32  `json:"chunk_num,omitempty"`
		FileId   int32  `json:"file_id"`
	}

	p := new(preUploadReq)

	if err := ctx.BodyParser(p); err != nil {
		return response.Resp400(ctx, nil)
	}
	//获取用户id
	id, ok := ctx.Locals("userID").(uint)
	fmt.Println(id)
	if !ok {
		return response.Resp400(ctx, "token expire time")
	}
	//创建用户信息
	createFileResp, err := a.fileRPCClient.CreateFile(context.Background(), &file.CreateFileReq{
		UserId:    int64(id),
		Path:      p.Path,
		FileName:  p.FileName,
		Ext:       p.Ext,
		Hash:      p.Hash,
		Size:      p.Size,
		IsPrivate: p.IsPrivate,
		StorageId: 0,
	})
	if err != nil {
		logrus.Error("create file rpc err:", err)
		return response.Resp400(ctx, nil)
	}

	//获取上传链接
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
		FileId:   createFileResp.FileId,
	})
}
