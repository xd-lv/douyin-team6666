package pack

import (
	"context"
	"main/dal/mysqldb"
)

type video struct {
	id            int64
	title         string
	author        *user
	playUrl       string
	coverUrl      string
	favoriteCount int64
	commentCount  int64
	isFavorite    bool
}

func WithVideo(videoID int64) *video {
	return &video{
		id: videoID,
	}
}

func (v *video) GetVideo(ctx context.Context) error {
	err := v.getMysqlVideo(ctx)
	if err != nil {
		return err
	}

	err = v.getRocksdbVideo(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (v *video) getIsFavorite(ctx context.Context) error {
	// TODO 获取是否点赞，需要jwt，参数可能需要添加
	return nil
}

func (v *video) getRocksdbVideo(ctx context.Context) error {
	// TODO 从rocksdb中查询获得vido点赞评论等字段数据
	return nil
}

func (v *video) getMysqlVideo(ctx context.Context) error {
	dbVideo, err := mysqldb.GetVideo(ctx, v.id)
	if err != nil {
		return err
	}

	v.title = dbVideo.Title
	v.coverUrl = dbVideo.CoverUrl
	v.playUrl = dbVideo.PlayUrl

	author := WithUser(dbVideo.Author)
	err = author.GetUser(ctx)
	if err != nil {
		return err
	}

	v.author = author

	return nil
}
