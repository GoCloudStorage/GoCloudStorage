package model

import (
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/snowflake"
	"gorm.io/gorm"
)

type StorageInfo struct {
	gorm.Model
	StorageId uint64 `json:"storage_id,omitempty"`
	IsRemote  bool   `json:"is_remote"`      // 是否远程存储
	Hash      string `json:"hash,omitempty"` // 文件hash
	Size      int    `json:"size,omitempty"` // 文件大小
	RealPath  string `json:"real_path"`      // 文件存储位置
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

// FirstByHash 通过hash查找文件
func (s *StorageInfo) FirstByHash(hash string) error {
	tx := pg.Client.Where("hash=?", hash).First(s)
	return tx.Error
}

// CreateStorage 创建存储
func (s *StorageInfo) CreateStorage() error {
	return pg.Client.Create(s).Error
}

func (s *StorageInfo) GetStorageByStorageId(storageID uint64) error {
	return pg.Client.Where("storage_id=?", storageID).First(s).Error
}

func (s *StorageInfo) UpdateStorage() error {
	return pg.Client.Model(s).Where("storage_id=?", s.StorageId).Updates(s).Error

}

func (s *StorageInfo) IsExistByKey(hash string) bool {
	err := s.FirstByHash(hash)
	if err == nil {
		return true
	}
	return false
}
