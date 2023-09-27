package model

import (
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"gorm.io/gorm"
)

type FileInfo struct {
	gorm.Model
	FileName string `json:"file_name,omitempty"`
	Ext      string `json:"ext,omitempty"`
	Hash     string `json:"hash,omitempty"`

	Size       int32  `json:"size,omitempty"`
	UploaderId uint   `json:"uploader_id,omitempty"`
	StorageId  uint64 `json:"storage_id,omitempty"`
}

func (f *FileInfo) Create() error {
	return pg.Client.Create(f).Error
}

func (f *FileInfo) FindOneByHash() error {
	return pg.Client.Model(f).Where("hash = ?", f.Hash).First(f).Error
}

func (f *FileInfo) FindOneByID(id int) error {
	return pg.Client.Model(f).Where("id = ?", id).First(f).Error
}

func (f *FileInfo) FindFileByUserIdAndFileInfo(userId uint, path string, fileName string, ext string) error {
	return pg.Client.Model(f).Where("uploader_id=? and path=? and file_name=? and ext=?", userId, path, fileName, ext).First(f).Error
}

func Init() {
	pg.Client.AutoMigrate(&FileInfo{})
}
