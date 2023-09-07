package minio

import (
	"work-space/pkg/storage"
)

type StorageEngine struct {
}

func (s *StorageEngine) Upload(request storage.UploadRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s *StorageEngine) GetTemporaryURL(fileMD5 string) string {
	//TODO implement me
	panic("implement me")
}

func (s *StorageEngine) GetPermanentURL(fileMD5 string) string {
	//TODO implement me
	panic("implement me")
}

func (s *StorageEngine) Init(config storage.InitConfig) {
	//TODO implement me
	panic("implement me")
}
