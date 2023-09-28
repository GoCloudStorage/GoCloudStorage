package oss

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
	"time"
)

type Minio struct {
	client     *minio.Client
	bucketname string
}

func NewMinio(endpoint string, username string, password, bucketname string) *Minio {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(username, password, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}

	return &Minio{client: client, bucketname: bucketname}
}

func (m *Minio) GetPreSignedDownloadURL(key string) (string, error) {
	u, err := m.client.PresignedGetObject(context.Background(), m.bucketname, key, time.Minute*15, nil)
	return u.String(), err
}

func (m *Minio) PutObject(key string, objectPath string, fileSize int64) (string, error) {
	var kk = fmt.Sprintf("object-key:%s", key)
	redis.SetLock(context.Background(), kk, time.Minute*30)
	defer redis.ReleaseLock(context.Background(), kk)
	file, err := os.OpenFile(objectPath, os.O_RDONLY, 0755)
	if err != nil {
		return "", fmt.Errorf("file open failed, err: %v", err)
	}
	info, err := m.client.PutObject(context.Background(), m.bucketname, key, file, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("put object failed, err: %v", err)
	}
	return info.ETag, nil
}

func (m *Minio) PutBucket() error {
	//TODO implement me
	panic("implement me")
}
