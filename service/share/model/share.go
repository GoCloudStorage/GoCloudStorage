package model

import (
	"gorm.io/gorm"
	"time"
)

type ShareURl struct {
	gorm.Model
	UserID int
	Url    string
	Hash   string
	Expire time.Time
}
