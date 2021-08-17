package main

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/Jsoneft/cos_service/constant"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	u, _ := url.Parse(constant.JarvisImageUploadURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  constant.SecretID,
			SecretKey: constant.SecretKey,
		},
	})
	// 多线程批量上传文件
	filesCh := make(chan string, 2)
	filePaths := []string{"test1", "test2", "test3"}
	var wg sync.WaitGroup
	threadpool := 2
	for i := 0; i < threadpool; i++ {
		wg.Add(1)
		go upload(&wg, c, filesCh)
	}
	for _, filePath := range filePaths {
		filesCh <- filePath
	}
	close(filesCh)
	wg.Wait()
}
