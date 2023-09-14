package local

import (
	"bytes"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage_engine"
	"os"
	"testing"
	"time"
)

func TestStorageEngineUpload(t *testing.T) {
	var s StorageEngine
	s.Init(storage_engine.InitConfig{
		Endpoint:        "./",
		AccessKeyID:     "",
		SecretAccessKey: "",
		UseSSL:          false,
		BucketName:      "test",
	})
	err := s.UploadChunk(storage_engine.UploadChunkRequest{
		FileMD5: "123456",
		Data:    bytes.NewReader([]byte("hello")),
		PartNum: 0,
	})
	if err != nil {
		t.Fatal("failed to upload chunk 1")
	}

	err = s.UploadChunk(storage_engine.UploadChunkRequest{
		FileMD5: "123456",
		Data:    bytes.NewReader([]byte(" world")),
		PartNum: 1,
	})
	if err != nil {
		t.Fatal("failed to upload chunk 2")
	}

	err = s.UploadChunk(storage_engine.UploadChunkRequest{
		FileMD5: "123456",
		Data:    bytes.NewReader([]byte(" cill")),
		PartNum: 2,
	})
	if err != nil {
		t.Fatal("failed to upload chunk 3")
	}

	err = s.MergeChunk("123456", 3, len([]byte("hello world cill")))
	if err != nil {
		t.Fatal("failed to merge chunk", err)
	}
}

func TestLocalStorageEngineGetURL(t *testing.T) {
	redis.Init("162.14.115.114:6379", "12345678", 0)
	var s StorageEngine
	curDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	s.Init(storage_engine.InitConfig{
		Endpoint:        curDir,
		AccessKeyID:     "",
		SecretAccessKey: "",
		UseSSL:          false,
		BucketName:      "test",
	})
	url, err := s.GenerateObjectURL("123456", time.Second*150)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}
