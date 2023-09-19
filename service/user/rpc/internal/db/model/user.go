package model

import (
	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
)

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

func Init() {
	pg.Client.AutoMigrate(User{})
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	// 生成盐值，并用盐值对密码进行哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 返回哈希后的密码作为字符串
	return string(hashedPassword), nil
}

// ComparePasswords 比较密码和已加密的密码是否匹配
func ComparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	logrus.Error(hashedPassword, password)
	return err
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
