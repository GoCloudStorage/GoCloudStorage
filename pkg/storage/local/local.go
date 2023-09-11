package local

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage"
	redis2 "github.com/redis/go-redis/v9"
	"path"
	"time"
)

var (
	FileAlreadyExist = fmt.Errorf("file already exist")
)

type StorageEngine struct {
	rootPath string
	uploader chunkUploader
}

func (s *StorageEngine) UploadChunk(request storage.UploadChunkRequest) error {
	dirPath := s.getFileDir(request.FileMD5)
	return s.uploader.saveChunk(dirPath, request.PartNum, request.Data)
}

func (s *StorageEngine) getFileDir(fileMD5 string) string {
	return path.Join(s.rootPath, fileMD5)
}

func (s *StorageEngine) GetObjectURL(fileMD5 string, expire time.Duration) string {
	filePath := path.Join(s.getFileDir(fileMD5), "data")
	cmd := redis.Client.Get(context.Background(), fileMD5)
	if cmd.Err() == redis2.Nil {
		redis.Client.SetEx(context.Background(), fileMD5, filePath, expire)
	} else if cmd.Err() != nil {
		return "Err"
	}
	return filePath
}

func (s *StorageEngine) Init(config storage.InitConfig) {
	s.rootPath = path.Join(config.Endpoint, config.BucketName)
}

func (s *StorageEngine) MergeChunk(fileMD5 string, partSize int, dataSize int) error {
	dirPath := s.getFileDir(fileMD5)
	return s.uploader.mergeChunk(dirPath, partSize, dataSize)
}
