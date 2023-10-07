package handler

import (
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (a *API) getAllFile(c *fiber.Ctx) error {
	var (
		uploadId uint
		fileDB   model.FileInfo
	)

	id, ok := c.Locals("userID").(uint)
	if !ok {
		logrus.Error("not have user id")
		return response.Resp400(c, nil)
	}
	uploadId = id

	fileInfos, err := fileDB.FindAllByUploaderID(uploadId)
	if err != nil {
		logrus.Errorf("find all file by uploadID failed, err: %v", err)
		return response.Resp500(c, nil)
	}

	return response.Resp200(c, fileInfos)
}
