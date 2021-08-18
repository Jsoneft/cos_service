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

	"github.com/Jsoneft/cos_service/constant/config"
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
	u, _ := url.Parse(config.Cos.UploadURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Timeout: config.Cos.COSTimeOut,
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Cos.SecretID,
			SecretKey: config.Cos.SecretKey,
		},
	})
	return &client{
		c,
		b,
	}
}

// Upload 幂等串行上传 如果已经上传则掠过.
func (c *client) Upload(ctx context.Context, files []string) ([]string, error) {
	var res []string
	for _, file := range files {
		key, fd, err := c.getFileKey(file)
		if err != nil {
			return nil, err
		}
		if c.isUploaded(ctx, key) {
			res = append(res, c.Client.Object.GetObjectURL(key).String())
			continue
		}
		_, err = c.Client.Object.Put(ctx, key, fd, nil)
		if err != nil {
			return nil, err
		}
		//res = append(res, fmt.Sprintf(config.Cos.ImgURLTemplate, MD5Str, fileSuffix))
		res = append(res, c.Client.Object.GetObjectURL(key).String())
	}
	return res, nil
}

// IsUploaded returns ture if the file is uploaded.
func (c *client) isUploaded(ctx context.Context, fileName string) bool {
	_, err := c.Client.Object.Head(ctx, fileName, nil)
	return !cos.IsNotFoundError(err)
}

// GetFileKey 返回文件的副本并拿到文件的唯一key.
func (c *client) getFileKey(file string) (string, *bytes.Reader, error) {
	fd, err := os.Open(file)
	defer fd.Close()
	if err != nil {
		return "", nil, err
	}
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return "", nil, err
	}
	fileSuffix := path.Ext(file)
	tmpFD := bytes.NewReader(content)
	content = append(content, []byte(fileSuffix)...)
	MD5Str := fmt.Sprintf("%x", md5.Sum(content))
	name := fmt.Sprintf("%s%s%s", config.Cos.COSRelativePath, MD5Str, fileSuffix)
	return name, tmpFD, nil
}
