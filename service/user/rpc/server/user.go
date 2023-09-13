package server

//客户端代码

import (
	"context"
	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
)

func UserRegister(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	resp, err = UserClient.UserRegister(ctx, req)
	if err != nil {
		return resp, errors.New("注册失败")
	}
	resp.St = true

	return resp, nil
}

func UserLogin(ctx context.Context, req *user.LoginRequest) (resp *user.UserDetailResponse, err error) {
	resp, err = UserClient.UserLogin(ctx, req)

	if err != nil {
		return nil, errors.New("登陆失败")
	}

	return resp, nil
}
