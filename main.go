package main

import (
	"context"
	"github.com/Jsoneft/cos_service/constant"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func upload(wg *sync.WaitGroup, c *cos.Client, files <-chan string) {
	defer wg.Done()
	for file := range files {
		name := "folder/" + file
		fd, err := os.Open(file)
		if err != nil {
			//ERROR
			continue
		}
		_, err = c.Object.Put(context.Background(), name, fd, nil)
		if err != nil {
			//ERROR
		}
	}
}

func main() {
	u, _ := url.Parse(constant.JarvisImageUploadURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
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
