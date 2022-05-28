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

// GetFollowCount 获取关注数量
func GetFollowCount(ctx context.Context, userId int64) (int64, error) {
	followKey := "follow_" + strconv.FormatInt(userId, 10)
	followCount, err := rdb.ZCard(ctx, followKey).Result()
	if err != nil {
		return 0, errors.New("redis server error")
	}
	return followCount, nil
}

// GetFollowerCount 获取粉丝数量
func GetFollowerCount(ctx context.Context, userId int64) (int64, error) {
	fansKey := "fans_" + strconv.FormatInt(userId, 10)
	followerCount, err := rdb.ZCard(ctx, fansKey).Result()
	if err != nil {
		return 0, errors.New("redis server error")
	}
	return followerCount, nil
}

// IsAFollowB 判断用户A是否关注了用户B
func IsAFollowB(ctx context.Context, userAId int64, userBId int64) (bool, error) {
	followKey := "follow_" + strconv.FormatInt(userAId, 10)
	userBIdStr := strconv.FormatInt(userBId, 10)
	if err := rdb.ZRank(ctx, followKey, userBIdStr).Err(); err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, errors.New("redis server error")
	}
	return true, nil
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
			return []pack.User{}, err
		}
		followCount, err := GetFollowCount(ctx, id)
		if err != nil {
			return []pack.User{}, err
		}
		folloerCount, err := GetFollowerCount(ctx, id)
		if err != nil {
			return []pack.User{}, err
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

// GetFollowerList 获取粉丝列表
func GetFollowerList(ctx context.Context, userId int64) ([]pack.User, error) {
	users := make([]pack.User, 0)
	fansKey := "fans_" + strconv.FormatInt(userId, 10)

	res, err := rdb.ZRevRange(ctx, fansKey, 0, -1).Result()
	if err != nil {
		return users, errors.New("redis server error")
	}

	for _, idStr := range res {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		username, err := mysqldb.GetUserNameByID(ctx, id)
		if err != nil {
			return []pack.User{}, err
		}
		followCount, err := GetFollowCount(ctx, id)
		if err != nil {
			return []pack.User{}, err
		}
		folloerCount, err := GetFollowerCount(ctx, id)
		if err != nil {
			return []pack.User{}, err
		}
		isFollow, err := IsAFollowB(ctx, userId, id)
		if err != nil {
			return []pack.User{}, err
		}
		users = append(users, pack.User{
			Id:           id,
			Name:         username,
			FollowCount:  followCount,
			FolloerCount: folloerCount,
			IsFollow:     isFollow,
		})
	}

	return users, nil
}
