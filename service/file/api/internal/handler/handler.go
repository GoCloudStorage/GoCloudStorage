package handler

import (
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/gofiber/fiber/v2"
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
		Token     token.UploadToken `json:"token,omitempty"`
		StorageId int64             `json:"storageId,omitempty"`
	}

	p := new(preUploadReq)

	if err := ctx.BodyParser(p); err != nil {
		return err
	}

	//验参

	num := p.Size/opt.Cfg.File.BlockSize + 1
	token, err := token.GenerateToken(p.Hash, num, p.Size)
	if err != nil {
		return err
	}

}

func PreDownload(ctx *fiber.Ctx) error {
	panic("not impl")
}

func GetAll(ctx *fiber.Ctx) error {
	panic("not impl")
}
