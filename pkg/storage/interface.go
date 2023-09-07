package storage

import "io"

type IStorage interface {
	Upload(request UploadRequest) error
	GetTemporaryURL(fileMD5 string) string
	GetPermanentURL(fileMD5 string) string
	Init(InitConfig)
}

func New() IStorage {

}

type UploadRequest struct {
	BucketName string
	Filename   string

	Data io.Reader
	Size int64 // data 大小, 单位字节

	PartNum int // 分段上传: 块号, -1表示不是分段上传
}

type InitConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}
