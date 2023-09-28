package logic

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/pb/file"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FileServer struct {
	file.UnimplementedFileServer
}

func (s *FileServer) CreateFile(ctx context.Context, in *file.CreateFileReq) (*file.CreateFileResp, error) {
	if in.Hash == "" {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}

	f := &model.FileInfo{
		FileName:   in.FileName,
		Path:       in.Path,
		Size:       in.Size,
		Ext:        in.Ext,
		UploaderId: uint(in.UserId),
		Hash:       in.Hash,
		StorageId:  uint64(in.StorageId),
		IsPrivate:  false,
	}
	err := f.Create()

	if err != nil {
		logrus.Error("CreateFile err:", err)
		return nil, errors.New(response.RPC_DB_ERROR)
	}

	return &file.CreateFileResp{}, nil
}

func (s *FileServer) FindFileByUserIdAndFileInfo(ctx context.Context, in *file.FindFileByUserIdAndFileInfoReq) (*file.FindFileByUserIdAndFileInfoResp, error) {
	if in.UserId == 0 {
		return nil, errors.New(response.RPC_PARAM_ERROR)
	}

	f := new(model.FileInfo)
	err := f.FindFileByUserIdAndFileInfo(uint(in.UserId), in.Path, in.FileName, in.Ext)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error("FindFileByUserIdAndFileInfo err:", err)
		return nil, errors.New(response.RPC_DB_ERROR)
	}

	//未找到
	if err != nil {
		return &file.FindFileByUserIdAndFileInfoResp{}, nil
	}

	return &file.FindFileByUserIdAndFileInfoResp{
		Hash:      f.Hash,
		StorageId: int64(f.StorageId),
	}, nil
}
