package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"work-space/opt"
	"work-space/service/storage/internal/model"
	"work-space/service/storage/internal/svc"
)

func Upload(ctx *fiber.Ctx) error {
	type UploadReq struct {
		FileName         string   `json:"file_name,omitempty"`
		Tags             []string `json:"tags,omitempty"`
		AccessPermission int8     `json:"access_permission,omitempty"`
	}

	type UploadResp struct {
		IsSuccess bool `json:"is_success"`
	}

	var (
		err error
		req UploadReq
		res UploadResp
	)

	if err = ctx.JSON(&req); err != nil {
		return err
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		var (
			info model.FileInfo
		)
		f, err := file.Open()
		if err != nil {
			return err
		}

		data := make([]byte, file.Size)
		_, err = f.Read(data)
		if err != nil {
			panic(err)
		}
		logrus.Println(file.Filename, file.Size, string(data))

		// 计算文件hash
		hash := sha256.New()
		info.Hash = fmt.Sprintf("%x", hash.Sum(data))

		if err = info.FindOneByHash(); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err = svc.Upload(opt.Cfg.Storage.BucketName, file.Filename, bytes.NewReader(data), file.Size); err != nil {
			return err
		}
	}
	res.IsSuccess = true
	data, err := json.Marshal(&res)
	if err != nil {
		return err
	}
	err = ctx.Send(data)
	if err != nil {
		logrus.Errorf("send response failed, err: %s", err)
	}
	return nil
}

func Download(ctx *fiber.Ctx) error {
	panic("not impl")
}

func GetAll(ctx *fiber.Ctx) error {
	panic("not impl")
}
