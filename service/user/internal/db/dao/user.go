package dao

import (
	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/internal/db/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao() *UserDao {
	return &UserDao{pg.Client}
}

// GetUserInfo 获取用户信息
func (dao *UserDao) GetUserInfo(req *user.RegisterRequest) (user *model.User, err error) {
	err = dao.Model(&model.User{}).Where("email=?", req.Email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return
}

// CreateUser 创建用户
func (dao *UserDao) CreateUser(req *user.RegisterRequest) (err error) {
	var u *model.User
	var count int64
	err = dao.AutoMigrate(&u)
	if err != nil {
		return err
	}
	dao.Model(&model.User{}).Where("email=?", req.Email).Count(&count)
	if count != 0 {
		return errors.New("user already exits")
	}
	_ = u.SetPassword(req.Password)
	u = &model.User{
		UserName:    req.UserName,
		PassWord:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Photo:       req.Photo,
		Status:      0,
		Permission:  1,
	}
	if err = dao.Model(&model.User{}).Create(&u).Error; err != nil {
		logrus.Error("create u error", err)
		return
	}
	return
}
