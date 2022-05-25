package pack

import (
	"context"
	"main/dal/mysqldb"
)

type Video struct {
	Id            int64
	Title         string
	Author        *User
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
}

func WithVideo(videoID int64) Video {
	return Video{
		Id: videoID,
	}
}

func (v *Video) GetVideo(ctx context.Context) error {
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

func (v *Video) getIsFavorite(ctx context.Context) error {
	// TODO 获取是否点赞，需要jwt，参数可能需要添加
	return nil
}

func (v *Video) getRocksdbVideo(ctx context.Context) error {
	// TODO 从rocksdb中查询获得vido点赞评论等字段数据
	return nil
}

func (v *Video) getMysqlVideo(ctx context.Context) error {
	dbVideo, err := mysqldb.GetVideo(ctx, v.Id)
	if err != nil {
		return err
	}

	v.Title = dbVideo.Title
	v.CoverUrl = dbVideo.CoverUrl
	v.PlayUrl = dbVideo.PlayUrl

	author := WithUser(dbVideo.Author)
	err = author.GetUser(ctx)
	if err != nil {
		return err
	}

	v.Author = author

	return nil
}
