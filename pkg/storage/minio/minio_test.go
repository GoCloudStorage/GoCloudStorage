package minio

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestMinio(t *testing.T) {
	Init("162.14.115.114:9000", "cill", "12345678", false)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	file := bytes.NewReader([]byte("text content"))
	err := Upload(ctx, UploadConfig{
		BucketName: "test-01",
		FileName:   "testfile.go",
		File:       file,
		Size:       int64(file.Len()),
	})
	if err != nil {
		t.Fatalf("failed to upload, err: %v", err)
	}
}
