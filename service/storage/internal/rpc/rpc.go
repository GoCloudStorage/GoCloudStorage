package rpc

import (
	"io"
	"work-space/pb/storage"
)

var server storage.StorageServer

func init() {
	server = new(StorageServer)
}

type StorageServer struct {
	storage.UnsafeStorageServer
}

func (s StorageServer) Upload(stream storage.Storage_UploadServer) error {
	var (
		data []byte
	)
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			panic("file read success, not upload minio")
		}
		if err != nil {
			return err
		}
		data = append(data, chunk.File...)
	}
}
