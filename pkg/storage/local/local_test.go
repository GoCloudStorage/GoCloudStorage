package local

import (
	"bytes"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage"
	"testing"
)

func TestStorageEngineUpload(t *testing.T) {
	var s StorageEngine
	s.Init(storage.InitConfig{
		Endpoint:        "./",
		AccessKeyID:     "",
		SecretAccessKey: "",
		UseSSL:          false,
		BucketName:      "test",
	})
	err := s.UploadChunk(storage.UploadChunkRequest{
		FileMD5: "123456",
		Data:    bytes.NewReader([]byte("hello")),
		PartNum: 0,
	})
	if err != nil {
		t.Fatal("failed to upload chunk 1")
	}

	err = s.UploadChunk(storage.UploadChunkRequest{
		FileMD5: "123456",
		Data:    bytes.NewReader([]byte(" world")),
		PartNum: 1,
	})
	if err != nil {
		t.Fatal("failed to upload chunk 2")
	}

	err = s.UploadChunk(storage.UploadChunkRequest{
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
