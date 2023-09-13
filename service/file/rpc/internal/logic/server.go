package logic

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/pb/file"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FileServer struct {
	file.UnimplementedFileServer
}

func (s *FileServer) CreateFile(ctx context.Context, in *file.CreateFileReq) (*file.CreateFileResp, error) {
	if in == nil {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}
	f := &model.FileInfo{
		FileName:   in.FileName,
		Path:       in.Path,
		Size:       in.Size,
		BlockSize:  in.BlockSize,
		Ext:        in.Ext,
		UploaderId: uint(in.UserId),
		Hash:       in.Hash,
		StorageId:  in.StorageId,
		IsPrivate:  false,
	}
	err := f.Create()
	if err != nil {
		logrus.Error("CreateFile err:", err)
		return nil, errors.New(response.RPC_DB_ERROR)
	}
	return nil, nil
}

func (s *FileServer) FindFileByUserIdAndFileInfo(ctx context.Context, in *file.FindFileByUserIdAndFileInfoReq) (*file.FindFileByUserIdAndFileInfoResp, error) {
	if in == nil {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}
	f := new(model.FileInfo)
	err := f.FindFileByUserIdAndFileInfo(uint(in.UserId), in.Path, in.FileName, in.Ext)
	if err != nil {
		logrus.Error("FindFileByUserIdAndFileInfo err:", err)
		return nil, errors.New(response.RPC_DB_ERROR)
	}
	return &file.FindFileByUserIdAndFileInfoResp{
		Hash:      f.Hash,
		StorageId: f.StorageId,
	}, nil
}
