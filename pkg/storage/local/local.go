package local

import (
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage"
	"path"
)

var (
	FileAlreadyExist = fmt.Errorf("file already exist")
)

type StorageEngine struct {
	rootPath string
	uploader chunkUploader
}

func (s *StorageEngine) UploadChunk(request storage.UploadChunkRequest) error {
	dirPath := s.getFileDir(request.BucketName, request.FileMD5)
	return s.uploader.saveChunk(dirPath, request.PartNum, request.Data)
}

func (s *StorageEngine) getFileDir(bucketName string, fileMD5 string) string {
	return path.Join(s.rootPath, bucketName, fileMD5)
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
	s.rootPath = path.Join(config.Endpoint, config.BucketName)
}

func (s *StorageEngine) MergeChunk(fileMD5 string, partSize int, dataSize int) error {

}
