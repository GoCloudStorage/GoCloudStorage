package model

import (
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"gorm.io/gorm"
)

type StorageInfo struct {
	gorm.Model
	StorageId  int64  `json:"storage_id,omitempty"`
	Hash       string `json:"hash,omitempty"`
	BlockSize  int    `json:"block_size,omitempty"`
	Size       int    `json:"size,omitempty"`        // 文件大小
	IsComplete bool   `json:"is_complete,omitempty"` // 文件完整性
	RealPath   string `json:"real_path,omitempty"`
}

func Init() {
	pg.Client.AutoMigrate(StorageInfo{})
}
