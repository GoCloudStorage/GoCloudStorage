package local

import (
	"work-space/pkg/storage"
)

type StorageEngine struct {
}

func (l *StorageEngine) Upload(request storage.UploadRequest) error {
	//TODO implement me
	panic("implement me")
}

func (l *StorageEngine) GetTemporaryURL(fileMD5 string) string {
	//TODO implement me
	panic("implement me")
}

func (l *StorageEngine) GetPermanentURL(fileMD5 string) string {
	//TODO implement me
	panic("implement me")
}

func (l *StorageEngine) Init(config storage.InitConfig) {
	//TODO implement me
	panic("implement me")
}
