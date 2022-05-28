package dal

import (
	"context"
	"fmt"
	"testing"
)

func Test_redisDB(t *testing.T) {
	ctx := context.Background()
	Init()
	err := RDB.Ping(ctx).Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	RDB.Set(ctx, "User:1:Name", "Tom", 0)
	get := RDB.Get(ctx, "User:1:Name")

	fmt.Println(get.Val())

}
