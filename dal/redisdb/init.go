package redisdb

import (
	"github.com/go-redis/redis/v8"
	"main/constants"
)

/*
	https://redis.uptrace.dev/guide/#why-go-redis
	使用该类型的Redis框架，该连接为介绍
*/

var RDB *redis.Client
var RDBComment *redis.Client

func init() {
	opt, err := redis.ParseURL(constants.RedisDefaultDSN)
	if err != nil {
		panic(err)
	}
	RDB = redis.NewClient(opt)

	opt, err = redis.ParseURL(constants.RedisDefaultDSN + constants.RedisCommentNumDSN)
	if err != nil {
		panic(err)
	}
	RDBComment = redis.NewClient(opt)
}
