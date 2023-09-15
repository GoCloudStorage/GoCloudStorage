package handler

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/xrpc"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (a *API) UserRegister(ctx *fiber.Ctx) (err error) {
	type User struct {
		gorm.Model
		UserName    string `json:"user_name"`
		PassWord    string `json:"pass_word"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Photo       string `json:"photo"`
		Status      uint64 `json:"status"`     // 1 在线 2 下线
		Permission  uint64 `json:"permission"` // 1 普通用户 2 管理员 3 超级管理员
	}

	req := &User{}
	if err = ctx.BodyParser(&req); err != nil {
		return err
	}

	userClient, err := xrpc.InitUserRPCClient(user.NewUserServiceClient)

	if err != nil {
		logrus.Error("InitRPCClient err", err)
		return err
	}
	_, err = userClient.NewSession().UserRegister(context.Background(), &user.RegisterRequest{
		UserName:    req.UserName,
		Password:    req.PassWord,
		Email:       req.Email,
		Photo:       req.Photo,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		logrus.Info(err)
		return err
	}
	return response.Resp200(ctx, "注册成功")
}

func (a *API) UserLogin(ctx *fiber.Ctx) (err error) {
	type request struct {
		PhoneNumber string `json:"phone_number"`
		PassWord    string `json:"pass_word"`
	}
	var req request

	err = ctx.BodyParser(&req)
	if err != nil {
		return err
	}
	var loginreq = user.LoginRequest{
		PhoneNumber: req.PhoneNumber,
		Password:    req.PassWord,
	}
	userClient, err := xrpc.InitUserRPCClient(user.NewUserServiceClient)

	if err != nil {
		logrus.Error("InitRPCClient err", err)
		return err
	}
	resp, err := userClient.NewSession().UserLogin(context.Background(), &loginreq)

	if err != nil {
		logrus.Info(err)
		return err
	}
	return response.Resp200(ctx, fiber.Map{
		"token": resp.Token,
		"msg":   "登录成功",
	})
}
