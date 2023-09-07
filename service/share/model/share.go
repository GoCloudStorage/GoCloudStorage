package model

import "gorm.io/gorm"

type ShareURl struct {
	gorm.Model
	UserID int
	Url    string
	Hash   string
}
