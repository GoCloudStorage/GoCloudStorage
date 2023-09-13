package response

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

const (
	MSG200 = "请求成功"
	MSG202 = "请求成功, 请稍后..."
	MSG400 = "请求参数错误"
	MSG401 = "登录已过期, 请重新登录"
	MSG403 = "请求权限不足"
	MSG404 = "请求资源未找到"
	MSG429 = "请求过于频繁, 请稍后再试"
	MSG500 = "服务器开小差了, 请稍后再试"
	MSG501 = "功能开发中, 尽情期待"
)

func handleEmptyMsg(status uint32, msg string) string {
	if msg == "" {
		switch status {
		case 200:
			msg = MSG200
		case 202:
			msg = MSG202
		case 400:
			msg = MSG400
		case 401:
			msg = MSG401
		case 403:
			msg = MSG403
		case 404:
			msg = MSG404
		case 429:
			msg = MSG429
		case 500:
			msg = MSG500
		case 501:
			msg = MSG501
		}
	}

	return msg
}

func Resp(c *fiber.Ctx, status uint32, msg string, data any) error {
	msg = handleEmptyMsg(status, msg)

	c.Set("SONAR-STATUS", strconv.Itoa(int(status)))

	if data == nil {
		return c.JSON(fiber.Map{"status": status, "msg": msg})
	}

	return c.JSON(fiber.Map{"status": status, "msg": msg, "data": data})
}

func Resp200(c *fiber.Ctx, data any, msgs ...string) error {
	msg := MSG200

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 200, msg, data)
}

func Resp202(c *fiber.Ctx, data any, msgs ...string) error {
	msg := MSG202

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 202, msg, data)
}

func Resp400(c *fiber.Ctx, data any, msgs ...string) error {
	msg := MSG400

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 400, msg, data)
}

func Resp401(c *fiber.Ctx, data any, msgs ...string) error {
	msg := MSG401

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 401, msg, data)
}

func Resp403(c *fiber.Ctx, data any, msgs ...string) error {
	msg := MSG403

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 403, msg, data)
}

func Resp429(c *fiber.Ctx, data any, msgs ...string) error {
	msg := MSG429

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 429, msg, data)
}

func Resp500(c *fiber.Ctx, data any, msgs ...string) error {
	msg := MSG500

	if len(msgs) > 0 && msgs[0] != "" {
		msg = fmt.Sprintf("%s: %s", msg, strings.Join(msgs, "; "))
	}

	return Resp(c, 500, msg, data)
}
