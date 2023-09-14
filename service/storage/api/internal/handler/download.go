package api

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/response"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (a *API) download(c *fiber.Ctx) error {
	// 获取range
	ranges, err := convertRange(c.Get("Content-Range"))
	start, end := ranges[0], ranges[1]
	if err != nil {
		return response.Resp400(c, nil, err.Error())
	}
	// 解析token
	//tokenVal := c.FormValue("token")
	//downloadToken, err := token.ParseDownloadToken(tokenVal)
	//if err != nil {
	//	return response.Resp400(c, nil, err.Error())
	//}

	// 获取文件路径
	key := c.Params("key")
	resp, err := a.storageRPCClient.GetRealPathByCode(context.Background(), &storage.GetRealPathByCodeReq{Code: key})
	if err != nil {
		return response.Resp500(c, nil, err.Error())
	}

	// 读取文件
	file, err := os.OpenFile(resp.Path, os.O_RDONLY, 0755)
	if err != nil {
		return response.Resp500(c, nil, "open file failed")
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return response.Resp500(c, nil, "get file info failed")
	}
	// 读取偏移
	size := min(opt.Cfg.File.BlockSize, min(fileInfo.Size()-start, end))
	body := make([]byte, size)
	file.Seek(ranges[0], 0)
	n, err := file.Read(body)
	size = int64(n)
	// 传输文件
	c.Set("Accept-Ranges", "bytes")
	c.Status(http.StatusPartialContent)
	c.Set("Range", fmt.Sprintf("bytes %d-%d/%d", start, start+size, fileInfo.Size()))
	return c.Send(body)
}

func convertRange(data string) ([]int64, error) {
	var res []int64
	tmp := strings.Split(data, "-")

	start, err := atoi64(tmp[0])
	if err != nil {
		return nil, fmt.Errorf("failed convert range")
	}
	res = append(res, start)
	end, err := atoi64(tmp[1])
	if err != nil {
		return nil, fmt.Errorf("failed convert range")
	}
	res = append(res, end)
	return res, nil
}

func atoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
