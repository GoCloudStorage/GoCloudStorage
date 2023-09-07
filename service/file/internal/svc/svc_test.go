package svc

import (
	"bytes"
	"testing"
	"work-space/opt"
	"work-space/pkg/storage/minio"
)

func TestUploadFile(t *testing.T) {
	opt.Cfg.Storage.BucketName = "test"
	minio.Init("162.14.115.114:9000", "cill", "12345678", false)
	str := "feafeafa"
	if err := Upload("test", "testname.txt", bytes.NewReader([]byte(str)), int64(len(str))); err != nil {
		t.Fatal("failed to upload storage ", err)
	}
}
