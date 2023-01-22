package utils

import (
	"context"
	"douyin/config"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

var (
	client    *cos.Client
	addr      string
	secretId  string
	secretKey string
)

// init 加载包时就获取client
func init() {
	cosConf := config.Conf.COS
	addr = cosConf.Addr
	secretId = cosConf.SecretId
	secretKey = cosConf.SecretKey
	client = GetClient()
}

// GetClient 返回cos client
func GetClient() *cos.Client {
	u, _ := url.Parse(addr)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			SecretKey: secretKey, // 用户的 secretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})
	return client
}

// UploadFile 上传文件
func UploadFile(key string, file *multipart.FileHeader) error {
	fd, err := file.Open()
	defer func(fd multipart.File) {
		err := fd.Close()
		if err != nil {
			log.Fatal("cannot open file", err)
		}
	}(fd)
	_, err = client.Object.Put(
		context.Background(), key, fd, nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetSignUrl 返回预签名Url
func GetSignUrl(key string) string {
	ctx := context.Background()
	presignedUrl, err := client.Object.GetPresignedURL(ctx, http.MethodGet, key, secretId, secretKey, time.Hour, nil)
	if err != nil {
		log.Fatal(err)
	}
	return presignedUrl.String()
}
