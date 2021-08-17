// Package client 封装了COS的业务逻辑。
package client

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"golang.org/x/sync/errgroup"

	"github.com/Jsoneft/cos_service/constant"
	"github.com/tencentyun/cos-go-sdk-v5"
)

//go:generate mockgen -destination=./mock/uploader_mock.go -package=mock . Client

// Client is the COS service implement.
type Client interface {
	ParallelUpload(ctx context.Context, files []string) error
}

type client struct {
	*cos.Client
	*cos.BaseURL
}

var New = func() Client {
	u, _ := url.Parse(constant.JarvisImageUploadURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Timeout: constant.COSTimeOut,
		Transport: &cos.AuthorizationTransport{
			SecretID:  constant.SecretID,
			SecretKey: constant.SecretKey,
		},
	})
	return &client{
		c,
		b,
	}
}

// ParallelUpload 并行上传 如某一个行为出错 立即停掉其他的上传并返回改次上传错误。
func (c *client) ParallelUpload(ctx context.Context, files []string) ([]string, error) {
	g, subCtx := errgroup.WithContext(ctx)
	for _, file := range files {
		file := file
		g.Go(func() error {
			fd, err := os.Open(file)
			if err != nil {
				return err
			}
			content, err := ioutil.ReadAll(fd)
			if err != nil {
				return err
			}
			content = append(content, []byte(file)...)
			encrypt := md5.New()
			MD5Str := hex.EncodeToString(encrypt.Sum(content))
			fileSuffix := path.Ext(file)
			name := fmt.Sprintf("%s%s.%s", constant.COSRelativePath, MD5Str, fileSuffix)
			rsp, err := c.Object.Put(subCtx, name, fd, nil)
			return err
		})
	}
	return g.Wait()
}
