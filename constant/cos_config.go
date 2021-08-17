package constant

import "time"

const (
	// JarvisImageUploadURL 个人COS服务访问域名
	JarvisImageUploadURL = "https://jarviszuo-tencent-img-1302316844.cos.ap-chengdu.myqcloud.com"

	// SecretID 访问服务密钥ID
	SecretID = "AKID74kAJCo9E3cuuf8KFtHOEn57sMaXFgqj"
	// SecretKey 访问服务密钥
	SecretKey = "dz4rNcaIirjKehTyDBtKVP2PbXBYeT1I"

	// COSTimeOut COS服务超时时间
	COSTimeOut = 100 * time.Second

	// COSRelativePath COS服务 上传路径
	COSRelativePath = "img/markdownImg/"
)
