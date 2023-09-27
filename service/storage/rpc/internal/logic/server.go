package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/oss"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/token"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/model"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

type StorageServer struct {
	HttpAddr string
	Oss      oss.IOss
	storage.UnimplementedStorageServer
}

func (s StorageServer) GetUploadURL(ctx context.Context, req *storage.GetUploadURLReq) (*storage.GetUploadURLResp, error) {
	panic("not impl")
}

func (s StorageServer) GetDownloadURL(ctx context.Context, req *storage.GetDownloadURLReq) (*storage.GetDownloadURLResp, error) {
	var (
		key         string
		filename    string
		ext         string
		expire      time.Duration
		storageInfo model.StorageInfo
	)
	// verify hash
	key = req.GetHash()
	if key == "" {
		return nil, errors.New("hash not is \"\"")
	}

	if err := storageInfo.FirstByHash(key); err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("[hash:%s] not in storage", key)
	}

	// verify filename
	filename = req.GetFilename()
	if filename == "" {
		filename = "default"
	}

	// verify ext
	ext = req.GetExt()
	if ext == "" {
		ext = "data"
	}

	// verify expire
	if t := req.GetExpire(); t == 0 {
		expire = time.Minute * 30
	} else {
		expire = time.Second * time.Duration(t)
	}

	// 通过OSS或者本地传输文件
	if storageInfo.IsRemote {
		// 生成OSS下载URL
		downloadURL, err := s.Oss.GetPreSignedDownloadURL(key)
		if err != nil {
			return nil, err
		}

		return &storage.GetDownloadURLResp{
			Url:       downloadURL,
			TotalSize: 0,
		}, nil
	} else {
		// 生成本地下载URL
		downloadToken, err := token.GenerateDownloadToken(storageInfo.StorageId, filename, ext, expire)
		if err != nil {
			return nil, fmt.Errorf("generate token fail, err: %v", err)
		}

		return &storage.GetDownloadURLResp{
			Url:       s.componentDownloadURL(downloadToken),
			TotalSize: 0,
		}, nil
	}
}

func (s StorageServer) componentDownloadURL(downloadToken string) string {
	result, err := url.JoinPath(s.HttpAddr, "/storage/download", downloadToken)
	if err != nil {
		logrus.Error(err)
	}
	return result
}
