package model

import (
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/snowflake"
	"gorm.io/gorm"
)

type StorageInfo struct {
	gorm.Model
	StorageId  uint64 `json:"storage_id,omitempty"`
	Hash       string `json:"hash,omitempty"`
	Size       int    `json:"size,omitempty"`        // 文件大小
	IsComplete bool   `json:"is_complete,omitempty"` // 文件完整性
	RealPath   string `json:"real_path,omitempty"`
}

func Init() {
	pg.Client.AutoMigrate(StorageInfo{})
}

func (s *StorageInfo) BeforeCreate(tx *gorm.DB) error {
	id, err := snowflake.GetID()
	if err != nil {
		return err
	}
	s.StorageId = id
	return nil
}

// FindStorageByHash 通过hash查找文件
func (s *StorageInfo) FindStorageByHash(hash string) error {
	tx := pg.Client.Where("hash=?", s.Hash).First(s)
	return tx.Error
}

// CreateStorage 创建存储
func (s *StorageInfo) CreateStorage() error {
	return pg.Client.Create(s).Error
}
