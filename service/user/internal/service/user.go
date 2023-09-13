package service

//服务端代码

import (
	"context"
	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/internal/db/dao"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/internal/db/model"
	"sync"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
	user.UnimplementedUserServiceServer
}

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (u *UserSrv) UserRegister(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
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

func (u *UserSrv) UserLogin(ctx context.Context, req *user.LoginRequest) (resp *user.UserDetailResponse, err error) {
	resp = new(user.UserDetailResponse)
	var usr model.User
	if err = dao.NewUserDao().Model(&model.User{}).Where("email=?", req.PhoneNumber).First(&usr).Error; err != nil {
		return resp, err
	}
	if ok := usr.CheckPassword(req.Password); !ok {
		return resp, errors.New("验证失败，密码错误")
	}

	//验证成功 生成并返回token
	userToken, err := token.GenToken(usr.ID)
	if err != nil {
		return resp, err
	}
	resp.Token = userToken
	return
}
