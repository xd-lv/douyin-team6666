package relationService

import (
	"context"
	"github.com/go-redis/redis/v8"
	"main/constants"
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
	GetRelation(ctx context.Context, user *pack.User) error
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

// getFollowKey 获取 Redis 关注键
func (rs *Impl) getFollowKey(userId int64) string {
	return constants.FollowKeyPrefix + strconv.FormatInt(userId, 10)
}

// getFollowKey 获取 Redis 粉丝键
func (rs *Impl) getFansKey(userId int64) string {
	return constants.FansKeyPrefix + strconv.FormatInt(userId, 10)
}

// Follow 关注操作
func (rs *Impl) Follow(ctx context.Context, userId int64, toUserId int64) error {
	timestamp := float64(time.Now().Unix())
	followKey := rs.getFollowKey(userId)
	fansKey := rs.getFansKey(toUserId)

	pipe := rdb.TxPipeline()
	pipe.ZAdd(ctx, followKey, &redis.Z{Score: timestamp, Member: toUserId})
	pipe.ZAdd(ctx, fansKey, &redis.Z{Score: timestamp, Member: userId})
	_, err := pipe.Exec(ctx)
	if err != nil {
		return constants.ErrRedisServer
	}

	return nil
}

// CancelFollow 取消关注操作
func (rs *Impl) CancelFollow(ctx context.Context, userId int64, toUserId int64) error {
	followKey := rs.getFollowKey(userId)
	fansKey := rs.getFansKey(toUserId)

	pipe := rdb.TxPipeline()
	pipe.ZRem(ctx, followKey, toUserId)
	pipe.ZRem(ctx, fansKey, userId)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return constants.ErrRedisServer
	}

	return nil
}

// getFollowCount 获取关注数量
func (rs *Impl) getFollowCount(ctx context.Context, userId int64) (int64, error) {
	followKey := rs.getFollowKey(userId)
	followCount, err := rdb.ZCard(ctx, followKey).Result()
	if err != nil {
		return 0, constants.ErrRedisServer
	}
	return followCount, nil
}

// getFollowerCount 获取粉丝数量
func (rs *Impl) getFollowerCount(ctx context.Context, userId int64) (int64, error) {
	fansKey := rs.getFansKey(userId)
	followerCount, err := rdb.ZCard(ctx, fansKey).Result()
	if err != nil {
		return 0, constants.ErrRedisServer
	}
	return followerCount, nil
}

// isAFollowB 判断用户A是否关注了用户B
func (rs *Impl) isAFollowB(ctx context.Context, userAId int64, userBId int64) (bool, error) {
	followKey := constants.FollowKeyPrefix + strconv.FormatInt(userAId, 10)
	userBIdStr := strconv.FormatInt(userBId, 10)
	if err := rdb.ZRank(ctx, followKey, userBIdStr).Err(); err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, constants.ErrRedisServer
	}
	return true, nil
}

// GetRelation 为 `用户`(登录用户|操作人) 提供 粉丝 及 关注 计数，字段 `IsFollow` 默认 false
func (rs *Impl) GetRelation(ctx context.Context, user *pack.User) error {
	var err error

	user.FollowCount, err = rs.getFollowCount(ctx, user.Id)
	if err != nil {
		return err
	}
	user.FollowerCount, err = rs.getFollowerCount(ctx, user.Id)
	if err != nil {
		return err
	}

	return nil
}

// GetRelationUser 获取 `用户`(登录用户|操作人) 的关系信息，包含关注数和粉丝数，字段 `IsFollow` 默认 false
func (rs *Impl) GetRelationUser(ctx context.Context, userId int64) (pack.User, error) {
	user := pack.User{Id: userId}

	username, err := mysqldb.GetUserNameByID(ctx, user.Id)
	if err != nil {
		return user, err
	}
	user.Name = username

	err = rs.GetRelation(ctx, &user)
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
	followKey := rs.getFollowKey(userId)

	res, err := rdb.ZRevRange(ctx, followKey, 0, -1).Result()
	if err != nil {
		return users, constants.ErrRedisServer
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
	fansKey := rs.getFansKey(userId)

	res, err := rdb.ZRevRange(ctx, fansKey, 0, -1).Result()
	if err != nil {
		return users, constants.ErrRedisServer
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
