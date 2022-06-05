package favoriteService

import (
	"context"
	"github.com/go-redis/redis/v8"
	"main/constants"
	"main/dal/redisdb"
	"strconv"
	"time"
)

var rdb = redisdb.RDB

type IFavoriteService interface {
	Like(ctx context.Context, userId int64, videoId int64) error
	Unlike(ctx context.Context, userId int64, videoId int64) error
	GetLikeList(ctx context.Context, userId int64) ([]int64, error)
	GetFavoriteCount(ctx context.Context, videoId int64) (int64, error)
	IsFavorite(ctx context.Context, userId int64, videoId int64) (bool, error)
}
type Impl struct {
}

func NewIFavoriteService() IFavoriteService {
	return &Impl{}
}

func (fs *Impl) Like(ctx context.Context, userId int64, videoId int64) error {
	timestamp := float64(time.Now().Unix())
	likeKey := "like_" + strconv.FormatInt(userId, 10)
	beLikedKey := "beLiked" + strconv.FormatInt(videoId, 10)

	pipe := rdb.TxPipeline()
	pipe.ZAdd(ctx, likeKey, &redis.Z{Score: timestamp, Member: videoId})
	pipe.SAdd(ctx, beLikedKey, userId)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return constants.ErrRedisServer
	}
	return nil
}

func (fs *Impl) Unlike(ctx context.Context, userId int64, videoId int64) error {
	likeKey := "like_" + strconv.FormatInt(userId, 10)
	beLikedKey := "beLiked" + strconv.FormatInt(videoId, 10)

	pipe := rdb.TxPipeline()
	pipe.ZRem(ctx, likeKey, videoId)
	pipe.SRem(ctx, beLikedKey, userId)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return constants.ErrRedisServer
	}
	return nil
}

func (fs *Impl) GetLikeList(ctx context.Context, userId int64) ([]int64, error) {
	videos := make([]int64, 0)
	likeKey := "like_" + strconv.FormatInt(userId, 10)
	res, err := rdb.ZRevRange(ctx, likeKey, 0, -1).Result()

	if err != nil {
		return videos, constants.ErrRedisServer
	}

	for _, idStr := range res {
		video, _ := strconv.ParseInt(idStr, 10, 64)
		videos = append(videos, video)
	}

	return videos, nil
}

func (fs *Impl) GetFavoriteCount(ctx context.Context, videoId int64) (int64, error) {
	beLikedKey := "beLiked" + strconv.FormatInt(videoId, 10)

	ret, err := rdb.SCard(ctx, beLikedKey).Result()
	if err != nil {
		return 0, constants.ErrRedisServer
	}

	return ret, err
}

func (fs *Impl) IsFavorite(ctx context.Context, userId int64, videoId int64) (bool, error) {
	beLikedKey := "beLiked" + strconv.FormatInt(videoId, 10)

	//toFind := strconv.FormatInt(userId, 10)
	ret, err := rdb.SIsMember(ctx, beLikedKey, userId).Result()

	if err != nil {
		return false, constants.ErrRedisServer
	}

	return ret, err
}
