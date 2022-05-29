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

type IRelationService interface {
	Follow(ctx context.Context, userId int64, toUserId int64) error
	CancelFollow(ctx context.Context, userId int64, toUserId int64) error
	GetRelationUser(ctx context.Context, userId int64) (pack.User, error)
	GetRelationAuthor(ctx context.Context, authorId int64, userId int64) (pack.User, error)
	GetFollowList(ctx context.Context, userId int64) ([]pack.User, error)
	GetFollowerList(ctx context.Context, userId int64) ([]pack.User, error)
}

type Impl struct {
}

func NewRelationService() IRelationService {
	return &Impl{}
}

// Follow 关注操作
func (rs *Impl) Follow(ctx context.Context, userId int64, toUserId int64) error {
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
func (rs *Impl) CancelFollow(ctx context.Context, userId int64, toUserId int64) error {
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

// getFollowCount 获取关注数量
func (rs *Impl) getFollowCount(ctx context.Context, userId int64) (int64, error) {
	followKey := "follow_" + strconv.FormatInt(userId, 10)
	followCount, err := rdb.ZCard(ctx, followKey).Result()
	if err != nil {
		return 0, errors.New("redis server error")
	}
	return followCount, nil
}

// getFollowerCount 获取粉丝数量
func (rs *Impl) getFollowerCount(ctx context.Context, userId int64) (int64, error) {
	fansKey := "fans_" + strconv.FormatInt(userId, 10)
	followerCount, err := rdb.ZCard(ctx, fansKey).Result()
	if err != nil {
		return 0, errors.New("redis server error")
	}
	return followerCount, nil
}

// isAFollowB 判断用户A是否关注了用户B
func (rs *Impl) isAFollowB(ctx context.Context, userAId int64, userBId int64) (bool, error) {
	followKey := "follow_" + strconv.FormatInt(userAId, 10)
	userBIdStr := strconv.FormatInt(userBId, 10)
	if err := rdb.ZRank(ctx, followKey, userBIdStr).Err(); err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, errors.New("redis server error")
	}
	return true, nil
}

// GetRelationUser 获取 `用户`(登录用户|操作人) 的关系信息，包含关注数和粉丝数，字段 `IsFollow` 默认 false
func (rs *Impl) GetRelationUser(ctx context.Context, userId int64) (pack.User, error) {
	user := pack.User{Id: userId}

	username, err := mysqldb.GetUserNameByID(ctx, user.Id)
	if err != nil {
		return user, err
	}
	user.Name = username

	user.FollowCount, err = rs.getFollowCount(ctx, user.Id)
	if err != nil {
		return user, err
	}
	user.FollowerCount, err = rs.getFollowerCount(ctx, user.Id)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetRelationAuthor 获取 `作者`(视频作者、评论作者等) 的关系信息，其中还包含 `用户` 对 `作者` 的关注情况
func (rs *Impl) GetRelationAuthor(ctx context.Context, authorId int64, userId int64) (pack.User, error) {
	user, err := rs.GetRelationUser(ctx, authorId)
	if err != nil {
		return user, err
	}

	user.IsFollow, err = rs.isAFollowB(ctx, userId, authorId)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetFollowList 获取关注列表
func (rs *Impl) GetFollowList(ctx context.Context, userId int64) ([]pack.User, error) {
	users := make([]pack.User, 0)
	followKey := "follow_" + strconv.FormatInt(userId, 10)

	res, err := rdb.ZRevRange(ctx, followKey, 0, -1).Result()
	if err != nil {
		return users, errors.New("redis server error")
	}

	for _, idStr := range res {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		user, err := rs.GetRelationAuthor(ctx, id, userId)
		if err != nil {
			return []pack.User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetFollowerList 获取粉丝列表
func (rs *Impl) GetFollowerList(ctx context.Context, userId int64) ([]pack.User, error) {
	users := make([]pack.User, 0)
	fansKey := "fans_" + strconv.FormatInt(userId, 10)

	res, err := rdb.ZRevRange(ctx, fansKey, 0, -1).Result()
	if err != nil {
		return users, errors.New("redis server error")
	}

	for _, idStr := range res {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		user, err := rs.GetRelationAuthor(ctx, id, userId)
		if err != nil {
			return []pack.User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}
