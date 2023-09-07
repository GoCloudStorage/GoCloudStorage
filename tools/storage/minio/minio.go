package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"io"
)

var (
	client *minio.Client
)

func Init(endpoint, accessKeyID, secretAccessKey string, useSSL bool) {
	var err error
	client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		logrus.Panicf("failed to init minio storage, err: %v", err)
	}
}

type UploadConfig struct {
	BucketName string
	FileName   string
	File       io.Reader
	Size       int64
}

func Upload(ctx context.Context, data UploadConfig) error {
	if exist, err := client.BucketExists(ctx, data.BucketName); err != nil {
		return err
	} else if !exist {
		err := client.MakeBucket(ctx, data.BucketName, minio.MakeBucketOptions{
			Region: "cn-north-1",
		})
		if err != nil {
			return err
		}
		policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement": [{"Effect":"Deny",
		"Principal": "*","Action": "s3:GetObject","Resource": "arn:aws:s3:::%s/*"}]}`, data.BucketName)

		err = client.SetBucketPolicy(ctx, data.BucketName, policy)
		if err != nil {
			return err
		}
	}

	_, err := client.PutObject(
		ctx,
		data.BucketName,
		data.FileName,
		data.File,
		data.Size,
		minio.PutObjectOptions{},
	)
	return err
}
