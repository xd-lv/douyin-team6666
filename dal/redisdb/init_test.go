package redisdb

import (
	"context"
	"fmt"
	"testing"
)

func Test_redisDB(t *testing.T) {
	ctx := context.Background()
	Init()
	err := RedisDB.Ping(ctx).Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	RedisDB.Set(ctx, "User:1:Name", "Tom", 0)
	get := RedisDB.Get(ctx, "User:1:Name")

	fmt.Println(get.Val())

}
