package api

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/gofiber/fiber/v2"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func Download(c *fiber.Ctx) error {
	// 解析token
	tokenVal := c.FormValue("token")
	downloadToken, err := token.ParseDownloadToken(tokenVal)
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}

	// 获取文件路径
	key := c.Params("key")
	cmd := redis.Client.Get(context.Background(), "storage:download:url"+key)
	if cmd.Err() == redis2.Nil {
		return response.Resp400(c, nil, response.MSG404)
	} else if cmd.Err() != nil {
		logrus.Errorf("redis client error, err %v", cmd.Err())
		return response.Resp500(c, nil)
	}
	filePath, err := cmd.Result()
	if err != nil {
		logrus.Errorf("failed to convert redis result, err: %v", err)
		return response.Resp500(c, nil)
	}
	logrus.Info(filePath)

	// 读取文件

	// 传输文件
	return response.Resp200(c, fmt.Sprintf("%s.%s", downloadToken.Filename, downloadToken.Ext), "success")
}
