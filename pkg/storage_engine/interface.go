package storage_engine

import (
	"io"
	"time"
)

var Client IStorage

type IStorage interface {
	Init(InitConfig)
	UploadChunk(request UploadChunkRequest) error
	MergeChunk(fileMD5 string, partSize int, dataSize int) error
	GenerateObjectURL(key string, expire time.Duration) (string, error)
	GetObjectURL(key string) (string, error)
}

type UploadChunkRequest struct {
	FileMD5 string
	Data    io.Reader
	PartNum int // 分段上传: 块号, -1表示不是分段上传
}

type InitConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}
