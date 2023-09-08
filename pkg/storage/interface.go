package storage

import "io"

type IStorage interface {
	Init(InitConfig)
	UploadChunk(request UploadChunkRequest) error
	MergeChunk(fileMD5 string, partSize int, dataSize int) error
	GetTemporaryURL(fileMD5 string) string
	GetPermanentURL(fileMD5 string) string
}

type UploadChunkRequest struct {
	BucketName string
	FileMD5    string

	Data io.Reader

	PartNum int // 分段上传: 块号, -1表示不是分段上传
}

type InitConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}
