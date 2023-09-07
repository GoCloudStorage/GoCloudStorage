package model

import (
	"gorm.io/gorm"
	"work-space/tools/db/pg"
)

type FileInfo struct {
	gorm.Model
	FileName  string `json:"file_name,omitempty"`
	Path      string `json:"path,omitempty"`
	Size      int64  `json:"size,omitempty"`
	BlockSize int64  `json:"blockSize"`
	Ext       string `json:"ext,omitempty"`
	Uploader  int    `json:"uploader,omitempty"`
	Hash      string `json:"hash,omitempty"`
}

func (f *FileInfo) FindOneByHash() error {
	return pg.Client.Model(f).Where("hash = ?", f.Hash).First(&f).Error
}

func Migreate() {
	pg.Client.AutoMigrate(&FileInfo{})
}
