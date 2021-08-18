package config

import "time"

type cos struct {
	// JarvisImageUploadURL 个人COS服务访问域名
	UploadURL string
	// SecretID 访问服务密钥ID
	SecretID string
	// SecretKey 访问服务密钥
	SecretKey string
	// COSTimeOut COS服务超时时间
	COSTimeOut time.Duration
	// COSRelativePath COS服务 上传路径
	COSRelativePath string
}

var Cos cos
