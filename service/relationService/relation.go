package relationService

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"main/dal/mysqldb"
	"main/dal/redisdb"
	"main/pack"
	"strconv"
	"time"
)

var rdb = redisdb.RDB

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
		return errors.New("redis server error")
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
		return errors.New("redis server error")
	}

	return nil
}

// GetFollowList 获取关注列表
func GetFollowList(ctx context.Context, userId int64) ([]pack.User, error) {
	users := make([]pack.User, 0)
	followKey := "follow_" + strconv.FormatInt(userId, 10)

	res, err := rdb.ZRevRange(ctx, followKey, 0, -1).Result()
	if err != nil {
		return users, errors.New("redis server error")
	}

	for _, idStr := range res {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		username, err := mysqldb.GetUserNameByID(ctx, id)
		if err != nil {
			return []pack.User{}, errors.New("mysql server error")
		}
		followCount, err := rdb.ZCard(ctx, "follow_"+idStr).Result()
		if err != nil {
			return []pack.User{}, errors.New("redis server error")
		}
		folloerCount, err := rdb.ZCard(ctx, "fans_"+idStr).Result()
		if err != nil {
			return []pack.User{}, errors.New("redis server error")
		}
		users = append(users, pack.User{
			Id:           id,
			Name:         username,
			FollowCount:  followCount,
			FolloerCount: folloerCount,
			IsFollow:     true,
		})
	}

	return users, nil
}
