package handler

import (
	"github.com/gofiber/fiber/v2"
	"work-space/pkg/token"
)

func PreUpload(ctx *fiber.Ctx) error {
	type UploadReq struct {
		FileName string `json:"file_name,omitempty" `
		Ext      string `json:"ext,omitempty"`
		Path     string `json:"path,omitempty"`
	}

	type UploadResp struct {
		Token     token.UploadToken `json:"token,omitempty"`
		UploadUrl string            `json:"upload_url,omitempty"`
	}

	panic("not impl")
}

func PreDownload(ctx *fiber.Ctx) error {
	panic("not impl")
}

func GetAll(ctx *fiber.Ctx) error {
	panic("not impl")
}
