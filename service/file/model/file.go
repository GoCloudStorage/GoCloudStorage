package model

import (
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"gorm.io/gorm"
)

type FileInfo struct {
	gorm.Model
	FileName   string `json:"file_name,omitempty"`
	Path       string `json:"path,omitempty"`
	Size       int64  `json:"size,omitempty"`
	BlockSize  int64  `json:"blockSize"`
	Ext        string `json:"ext,omitempty"`
	UploaderId int64  `json:"uploader,omitempty"`
	Hash       string `json:"hash,omitempty"`
	StorageId  int64  `json:"storage_id,omitempty"`
}

func (f *FileInfo) FindOneByHash() error {
	return pg.Client.Model(f).Where("hash = ?", f.Hash).First(&f).Error
}

func Init() {
	pg.Client.AutoMigrate(&FileInfo{})
}
