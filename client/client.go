// Package client 封装了COS的业务逻辑。
package client

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/Jsoneft/cos_service/constant"
	"github.com/tencentyun/cos-go-sdk-v5"
)

//go:generate mockgen -destination=./mock/uploader_mock.go -package=mock . Client

// Client is the COS service implement.
type Client interface {
	Upload(ctx context.Context, files []string) ([]string, error)
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

// Upload 串行上传.
func (c *client) Upload(ctx context.Context, files []string) ([]string, error) {
	var res []string
	for _, file := range files {
		fd, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		content, err := ioutil.ReadAll(fd)
		tmpFD := bytes.NewReader(content)
		if err != nil {
			return nil, err
		}
		fileSuffix := path.Ext(file)
		content = append(content, []byte(fileSuffix)...)
		MD5Str := fmt.Sprintf("%x", md5.Sum(content))
		name := fmt.Sprintf("%s%s%s", constant.COSRelativePath, MD5Str, fileSuffix)
		_, err = c.Client.Object.Put(ctx, name, tmpFD, nil)
		if err != nil {
			return nil, err
		}
		res = append(res, fmt.Sprintf(constant.ImgURLTemplate, MD5Str, fileSuffix))
	}
	return res, nil
}
