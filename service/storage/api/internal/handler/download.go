package handler

import (
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (a *API) download(c *fiber.Ctx) error {
	var (
		downloadToken *token.DownloadToken
		err           error
		storageInfo   model.StorageInfo
	)
	if t := c.Params("token"); t == "" {
		return response.Resp400(c, nil)
	} else {
		downloadToken, err = token.ParseDownloadToken(t)
		if err != nil {
			logrus.Error("verify token failed: ", err)
			return response.Resp400(c, nil)
		}
	}

	if err := storageInfo.GetStorageByStorageId(downloadToken.StorageID); err != nil {
		logrus.Error(err)
		return response.Resp400(c, nil)
	}
	// 设置响应头
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fmt.Sprintf("%s.%s", downloadToken.Filename, downloadToken.Ext)))
	c.Set("Content-Type", "application/"+downloadToken.Ext)

	return c.SendFile(storageInfo.RealPath)
}

func convertRange(data string) ([]int64, error) {
	var res []int64
	tmp := strings.Split(data, "=")
	tmp = strings.Split(tmp[1], "-")
	start, err := atoi64(tmp[0])
	if err != nil {
		return nil, fmt.Errorf("failed convert range")
	}
	res = append(res, start)
	end, err := atoi64(tmp[1])
	if err != nil {
		return nil, fmt.Errorf("failed convert range")
	}
	res = append(res, end)
	return res, nil
}

func atoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
