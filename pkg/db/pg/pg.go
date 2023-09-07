package pg

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Client *gorm.DB

func Init(host, user, password, dbname, port string) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	Client, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		logrus.Panicf("failed to connect pg, err: %v", err)
	}
}
