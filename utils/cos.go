package utils

import (
	"context"
	"douyin/common"
	"douyin/config"
	"douyin/models"
	"github.com/gomodule/redigo/redis"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"log"
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
func UploadFile(key string, file io.Reader) error {
	_, err := client.Object.Put(
		context.Background(), key, file, nil,
	)
	return err
}

// GetSignUrl 返回预签名Url
func GetSignUrl(key string) string {
	ctx := context.Background()
	redisConn := models.GetRedis()
	defer redisConn.Close()
	redisKey := common.RedisPrefixCos + common.RedisKeySplit + key
	signedUrl, err := redis.String(redisConn.Do("GET", redisKey))
	if err == nil && signedUrl != "" {
		return signedUrl
	}
	presignedUrl, err := client.Object.GetPresignedURL(ctx, http.MethodGet, key, secretId, secretKey, time.Hour, nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	err = redisConn.Send("SET", redisKey, presignedUrl.String())
	if err != nil {
		log.Println(err)
		return ""
	}
	err = redisConn.Send("EXPIRE", redisKey, 60*60-1)
	if err != nil {
		log.Println(err)
		return ""
	}
	_, err = redisConn.Do("")
	if err != nil {
		log.Println(err)
		return ""
	}
	redisConn.Close()
	return presignedUrl.String()
}
