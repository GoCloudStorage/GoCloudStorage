package model

import (
	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"gorm.io/gorm"
)

type FileInfo struct {
	gorm.Model
	UploaderId uint   `json:"uploader_id,omitempty"`
	FileName   string `json:"file_name,omitempty"`
	Path       string `json:"path,omitempty"`
	Ext        string `json:"ext,omitempty"`
	Hash       string `json:"hash,omitempty"`
	Size       int32  `json:"size,omitempty"`
	StorageId  uint64 `json:"storage_id,omitempty"`
	IsPrivate  bool   `json:"is_private"`
}

func (f *FileInfo) Create() error {
	return pg.Client.Create(f).Error
}

func (f *FileInfo) FindOneByHash() (error, bool) {
	err := pg.Client.Model(f).Where("hash = ?", f.Hash).First(f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false
	}
	return err, true
}

func (f *FileInfo) FindOneByID(id uint) error {
	return pg.Client.Model(f).Where("id = ?", id).First(f).Error
}

func (f *FileInfo) FindAllByUploaderID(id uint) (fileInfos []FileInfo, err error) {
	err = pg.Client.Model(f).Where("uploader_id = ?", id).Find(&fileInfos).Error
	return
}

func (f *FileInfo) UpdateFile() error {
	return pg.Client.Model(f).Where("id=?", f.ID).Updates(f).Error
}

func (f *FileInfo) FindFileByUserIdAndFileInfo(userId uint, path string, fileName string, ext string) error {
	return pg.Client.Model(f).Where("uploader_id=? and path=? and file_name=? and ext=?", userId, path, fileName, ext).First(f).Error
}

func Init() {
	pg.Client.AutoMigrate(&FileInfo{})
}
