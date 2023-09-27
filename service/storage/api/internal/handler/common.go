package handler

import (
	"errors"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/api/internal/logic"
	"github.com/gofiber/fiber/v2"

	"strconv"
	"strings"
)

func convertContentRange(d string) (res logic.ContentRange, err error) {
	t := strings.Split(d, " ")
	if len(t) != 2 {
		return res, fmt.Errorf("Content-Range format incorrect")
	}
	t = strings.Split(t[1], "-")
	if len(t) != 2 {
		return res, fmt.Errorf("Content-Range format incorrect")
	}
	res.Start, err = strconv.Atoi(t[0])
	if err != nil {
		return res, fmt.Errorf("start convert to int64 incorrect, err: %v", err)
	}
	t = strings.Split(t[1], "/")
	if len(t) != 2 {
		return res, fmt.Errorf("Content-Range format incorrect")
	}
	res.End, err = strconv.Atoi(t[0])
	if err != nil {
		return res, fmt.Errorf("end convert to int64 incorrect, err: %v", err)
	}
	res.Total, err = strconv.Atoi(t[1])
	if err != nil {
		return res, fmt.Errorf("total convert to int64 incorrect, err: %v", err)
	}
	return res, nil
}

func parasUploadHeader(c *fiber.Ctx) (*logic.UploadReq, error) {
	req, err := parasUploadCommonHeader(c)
	if err != nil {
		return nil, err
	}
	// 获取chunk number
	cn := c.Get("OSS-Chunk-Number", "")
	if cn == "" {
		return nil, errors.New("OSS-Chunk-Number not is nil")
	}
	chunkNumber, err := strconv.Atoi(cn)
	if err != nil {
		return nil, fmt.Errorf("OSS-Chunk-Number convert int fail, err: %v", err)
	}
	req.ChunkNumber = chunkNumber
	return req, nil
}

func parasUploadCommonHeader(c *fiber.Ctx) (*logic.UploadReq, error) {
	// 获取Content-Range
	r := c.Get("Content-Range", "nil")
	if r == "nil" { // 没有断点续传, 覆盖上传
		return nil, errors.New("Content-Range not set")
	}
	cr, err := convertContentRange(r)
	if err != nil {
		return nil, err
	}

	// 获取object total Size
	ocl := c.Get("OSS-Content-Length", "")
	if ocl == "" {
		return nil, errors.New("OSS-Content-Length not is nil")
	}
	total, err := strconv.Atoi(ocl)
	if err != nil {
		return nil, fmt.Errorf("convert OSS-Content-Length fail, err: %v", err)
	}

	// 获取object md5
	md5 := c.Get("OSS-MD5", "nil")
	if md5 == "nil" {
		return nil, errors.New("md5 form not found")
	}

	// 获取key
	key := c.Get("OSS-Key", "")
	if key == "" {
		return nil, errors.New("key form not is nil")
	}

	return &logic.UploadReq{
		ContentRange: cr,
		TotalSize:    total,
		Key:          key,
		ChunkNumber:  0,
	}, nil
}
