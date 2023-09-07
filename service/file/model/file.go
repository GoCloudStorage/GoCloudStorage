package model

import "gorm.io/gorm"

type StorageInfo struct {
	gorm.Model
	Hash       string
	BlockSize  int
	Size       int  // 文件大小
	IsComplete bool // 文件完整性
	RealPath   string
}
