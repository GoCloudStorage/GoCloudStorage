package http

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/rpc/server"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

func UserRegister(ctx *fiber.Ctx) (err error) {
	var req user.RegisterRequest

	if err = ctx.BodyParser(&req); err != nil {
		return response.Resp400(ctx, "绑定参数错误")
	}

	_, err = server.UserRegister(context.Background(), &req)
	if err != nil {
		logrus.Info(err)
		return response.Resp500(ctx, "注册失败")
	}

	return response.Resp200(ctx, "注册成功")
}

func UserLogin(ctx *fiber.Ctx) (err error) {
	var req user.LoginRequest

	if err = ctx.BodyParser(&req); err != nil {
		return response.Resp400(ctx, "绑定参数错误")
	}

	resp, err := server.UserLogin(context.Background(), &req)
	if err != nil {
		return response.Resp500(ctx, "登陆失败")
	}

	return response.Resp200(ctx, fiber.Map{
		"token": resp.Token,
		"msg":   "登录成功",
	})
}
