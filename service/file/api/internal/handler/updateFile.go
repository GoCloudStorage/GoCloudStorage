package handler

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/pb/file"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (a *API) updateFile(ctx *fiber.Ctx) error {
	type updateFile struct {
		FileId int64 `json:"file_id,required" form:"file_id,required"`

		FileName  string `json:"file_name" form:"file_name" `
		Ext       string `json:"ext" form:"ext"`
		Path      string `json:"path" form:"path"`
		IsPrivate bool   `json:"is_private" form:"is_private"`
		StorageId int64  `json:"storage_id" form:"storage_id"`
	}

	info := new(updateFile)
	err := ctx.BodyParser(info)
	if err != nil {
		return response.Resp400(ctx, nil)
	}

	//获取用户id
	id, ok := ctx.Locals("userID").(int64)
	if !ok {
		return response.Resp400(ctx, "token expire time")
	}

	_, err = a.fileRPCClient.UpdateFile(context.Background(), &file.UpdateFileReq{
		FileId:    info.FileId,
		UserId:    id,
		StorageId: info.StorageId,
		FileName:  info.FileName,
		Ext:       info.Ext,
		Path:      info.Path,
		IsPrivate: info.IsPrivate,
	})
	if err != nil {
		logrus.Error("update file err:", err)
		return response.Resp400(ctx, nil)
	}
	return response.Resp200(ctx, nil)
}
