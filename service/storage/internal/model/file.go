package model

import (
	"gorm.io/gorm"
	"work-space/tools/db/pg"
)

type FileInfo struct {
	gorm.Model
	FileName string `json:"file_name,omitempty"`
	Path     string `json:"path,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Ext      string `json:"ext,omitempty"`
	//Tags             []string `json:"tags,omitempty"`
	Uploader         int    `json:"uploader,omitempty"`
	AccessPermission int    `json:"access_permission,omitempty"`
	Hash             string `json:"hash,omitempty"`
	StorageLocation  string `json:"storage_location,omitempty"`
}

func (f *FileInfo) FindOneByHash() error {
	return pg.Client.Model(f).Where("hash = ?", f.Hash).First(&f).Error
}

func Migrator() {
	pg.Client.AutoMigrate(&FileInfo{})
}
