package handler

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/model"
	"github.com/gofiber/fiber/v2"
)

func (a *API) preDownload(c *fiber.Ctx) error {
	var (
		fileInfo model.FileInfo
	)
	fid, err := c.ParamsInt("id")
	if err != nil {
		return response.Resp400(c, nil, response.MSG400)
	}

	// 查找文件
	if err = fileInfo.FindOneByID(uint(fid)); err != nil {
		return response.Resp400(c, nil, "文件不存在", err.Error())
	}

	// 校验访问权限
	id, ok := c.Locals("userID").(uint)
	if !ok {
		return response.Resp400(c, nil, "not have user id")
	}
	var userid = id

	if fileInfo.IsPrivate && userid != fileInfo.UploaderId {
		return response.Resp400(c, nil, "没有访问权限")
	}

	// 生成下载链接, 调用storage server
	req := storage.GetDownloadURLReq{
		Hash:     fileInfo.Hash,
		Expire:   30,
		Filename: fileInfo.FileName,
		Ext:      fileInfo.Ext,
	}
	resp, err := a.storageRPCClient.GetDownloadURL(context.Background(), &req)
	if err != nil {
		return response.Resp500(c, nil, err.Error())
	}
	url := resp.GetUrl()

	// 返回数据
	type preDownloadResp struct {
		URL string `json:"url,omitempty"`
	}

	return response.Resp200(c, preDownloadResp{URL: url}, "success")
}
