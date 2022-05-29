package redisdb

import (
	"context"
	"fmt"
	"testing"
)

func Test_redisDB(t *testing.T) {
	ctx := context.Background()
	err := RDB.Ping(ctx).Err()
	if err != nil {
		fmt.Println("connect redis failed:", err)
		return
	}
	fmt.Println("redis connects successfully")
}
