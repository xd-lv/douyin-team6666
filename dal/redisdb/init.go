package redisdb

import (
	"github.com/go-redis/redis/v8"
	"main/constants"
)

/*
	https://redis.uptrace.dev/guide/#why-go-redis
	使用该类型的Redis框架，该连接为介绍
*/

var RedisDB *redis.Client

func Init() {
	opt, err := redis.ParseURL(constants.RedisDefaultDSN)
	if err != nil {
		panic(err)
	}
	RedisDB = redis.NewClient(opt)

}
