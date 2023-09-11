package util

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"strings"
)

var Msg = map[int]string{
	200: "请求成功",
	400: "请求参数错误",
	401: "登陆已过期，请重新登录",
	403: "请求权限不足",
	500: "服务器错误",
}

func Resp(c *fiber.Ctx, status uint32, msg string, data any) error {

	c.Set("X-Status", cast.ToString(status))

	if data == nil {
		return c.JSON(fiber.Map{"status": status, "msg": msg})
	}

	return c.JSON(fiber.Map{"status": status, "msg": msg, "data": data})
}

func Resp200(c *fiber.Ctx, data any, msgs ...string) error {
	msg := Msg[200]

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 200, msg, data)
}

func Resp400(c *fiber.Ctx, data any, msgs ...string) error {
	msg := Msg[400]

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 400, msg, data)
}

func Resp401(c *fiber.Ctx, data any, msgs ...string) error {
	msg := Msg[401]

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 401, msg, data)
}

func Resp403(c *fiber.Ctx, data any, msgs ...string) error {
	msg := Msg[403]

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 403, msg, data)
}

func Resp500(c *fiber.Ctx, data any, msgs ...string) error {
	msg := Msg[500]

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 500, msg, data)
}
