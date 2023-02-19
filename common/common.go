package common

import (
	"github.com/hertz-contrib/jwt"
	"strconv"
)

// 用来存放一些共用变量，防止循环引用
const (
	RedisKeySplit         = ":"
	RedisPrefixFavorVideo = "favorite:video"
	RedisPrefixRelation   = "relation:user"
	RedisFollowerField    = "follower"
	RedisFolloweeField    = "followee"
	RedisPrefixCos        = "cos:presign"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
)

func GetRedisRelationField(userId int) string {
	return RedisPrefixRelation + RedisKeySplit + strconv.Itoa(userId)
}
