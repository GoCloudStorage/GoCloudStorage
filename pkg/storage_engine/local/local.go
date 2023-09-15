package local

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage_engine"
<<<<<<< HEAD
	redis2 "github.com/redis/go-redis/v9"
	"log"
=======
>>>>>>> f5f05860dc07a675e4e61571dfb88bb9103fede2
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

<<<<<<< HEAD
=======
func New(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) storage_engine.IStorage {
	client := &StorageEngine{}
	client.Init(storage_engine.InitConfig{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		UseSSL:          useSSL,
		BucketName:      bucketName,
	})
	return client
}

>>>>>>> f5f05860dc07a675e4e61571dfb88bb9103fede2
func (s *StorageEngine) UploadChunk(request storage_engine.UploadChunkRequest) error {
	dirPath := s.getFileDir(request.FileMD5)
	return s.uploader.saveChunk(dirPath, request.PartNum, request.Data)
}

func (s *StorageEngine) getFileDir(fileMD5 string) string {
	return path.Join(s.rootPath, fileMD5)
}

// GenerateObjectURL 获取文件存储位置
<<<<<<< HEAD
func (s *StorageEngine) GenerateObjectURL(key, fileMD5 string, expire time.Duration) string {
	filePath := path.Join(s.getFileDir(fileMD5), "data")
	cmd := redis.Client.Get(context.Background(), key)
	if cmd.Err() == redis2.Nil {
		redis.Client.SetEx(context.Background(), key, filePath, expire)
	} else if cmd.Err() != nil {
		return "Err"
	}
	log.Println(cmd.Result())
	return filePath
=======
func (s *StorageEngine) GenerateObjectURL(fileMD5 string, expire time.Duration) (string, error) {
	filePath := path.Join(s.getFileDir(fileMD5), "data")
	return filePath, nil
>>>>>>> f5f05860dc07a675e4e61571dfb88bb9103fede2
}

func (s *StorageEngine) Init(config storage_engine.InitConfig) {
	s.rootPath = path.Join(config.Endpoint, config.BucketName)
}

func (s *StorageEngine) MergeChunk(fileMD5 string, partSize int, dataSize int) error {
	dirPath := s.getFileDir(fileMD5)
	return s.uploader.mergeChunk(dirPath, partSize, dataSize)
}

func (s *StorageEngine) GetObjectURL(key string) (string, error) {
	cmd := redis.Client.Get(context.Background(), key)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Result()
}
