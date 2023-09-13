package model

import (
	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
)

type User struct {
	gorm.Model
	UserName    string `json:"user_name"`
	PassWord    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Photo       string `json:"photo"`
	Status      uint64 `json:"status"`     // 1 在线 2 下线
	Permission  uint64 `json:"permission"` // 1 普通用户 2 管理员 3 超级管理员
}

func Init() {
	pg.Client.AutoMigrate(User{})
}

// SetPassword 加密密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user.PassWord = string(bytes)
	return nil
}

// CheckPassword 检验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(password))
	return err == nil
}

// UserMsgIsOk 判断用户信息是否符合要求
func UserMsgIsOk(req *user.RegisterRequest) (bool, error) {
	if len(req.UserName) < 3 || len(req.UserName) > 6 {
		return false, errors.New("username的长度应在3-6")
	}

	if len(req.Password) < 8 || len(req.Password) > 14 {
		return false, errors.New("用户密码应在8和14位")
	}

	pattern := `^1[3-9]\d{9}$`
	regex := regexp.MustCompile(pattern)
	//return regex.MatchString(phoneNumber)
	if !regex.MatchString(req.PhoneNumber) {
		return false, errors.New("您的手机号不符合规范")
	}
	return true, nil
}
