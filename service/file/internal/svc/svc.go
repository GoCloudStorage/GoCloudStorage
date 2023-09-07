package svc

import (
	"context"
	"io"
	"work-space/pkg/storage/minio"
)

func Upload(bucketName, filename string, data io.Reader, size int64) error {
	err := minio.Upload(context.Background(), minio.UploadConfig{
		BucketName: bucketName,
		FileName:   filename,
		File:       data,
		Size:       size,
	})
	if err != nil {
		return err
	}
	return nil
}
