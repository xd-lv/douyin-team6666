package relationService

import (
	"context"
	"github.com/go-redis/redis/v8"
	"main/dal/redisdb"
	"strconv"
	"time"
)

var rdb = dal.RDB

// Follow 关注操作
func Follow(ctx context.Context, userId int64, toUserId int64) error {
	timestamp := float64(time.Now().Unix())
	followKey := "follow_" + strconv.FormatInt(userId, 10)
	fansKey := "fans_" + strconv.FormatInt(toUserId, 10)

	pipe := rdb.TxPipeline()
	pipe.ZAdd(ctx, followKey, &redis.Z{Score: timestamp, Member: toUserId})
	pipe.ZAdd(ctx, fansKey, &redis.Z{Score: timestamp, Member: userId})
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// CancelFollow 取消关注操作
func CancelFollow(ctx context.Context, userId int64, toUserId int64) error {
	followKey := "follow_" + strconv.FormatInt(userId, 10)
	fansKey := "fans_" + strconv.FormatInt(toUserId, 10)

	pipe := rdb.TxPipeline()
	pipe.ZRem(ctx, followKey, toUserId)
	pipe.ZRem(ctx, fansKey, userId)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
