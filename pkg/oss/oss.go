package oss

type IOss interface {
	// GetPreSignedDownloadURL 获取下载预签名URL
	GetPreSignedDownloadURL(key string) (string, error)
	// PutBucket 创建存储桶
	PutBucket() error
}
