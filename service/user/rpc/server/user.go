package server

//服务端代码

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/rpc/internal/db/dao"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/rpc/internal/db/model"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
}

func (u *UserServiceServer) UserRegister(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	resp = new(user.RegisterResponse)
	resp.St = true
	if ok, err := model.UserMsgIsOk(req); !ok {
		resp.St = false
		return resp, err
	}

	err = dao.NewUserDao().CreateUser(req)
	if err != nil {
		resp.St = false
		return
	}
	return
}

func (u *UserServiceServer) UserLogin(ctx context.Context, req *user.LoginRequest) (resp *user.UserDetailResponse, err error) {
	resp = new(user.UserDetailResponse)
	var usr model.User
	if err = dao.NewUserDao().Model(&model.User{}).Where("phone_number=?", req.PhoneNumber).First(&usr).Error; err != nil {
		return resp, err
	}
	if err = model.ComparePasswords(usr.PassWord, req.Password); err != nil {
		return resp, errors.New(fmt.Sprintf("密码检验错误 %v", err))
	}

	//验证成功 生成并返回token
	userToken, err := token.GenToken(usr.ID)
	if err != nil {
		return resp, err
	}
	resp.Token = userToken
	return
}
