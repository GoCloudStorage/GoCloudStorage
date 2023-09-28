package handler

import (
	"github.com/GoCloudstorage/GoCloudstorage/service/file/model"
	"github.com/gofiber/fiber/v2"
)

func (a *API) getAllFile(c *fiber.Ctx) {
	var (
		uploadId int
		fileDB   model.FileInfo
	)

	fileDB.FindAllByUploaderID(uploadId)
}
