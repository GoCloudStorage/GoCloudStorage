package main

import "net/http"

func main() {
	http.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		// 验证token, 并确定上传分块
		// 上传
	})
}
