package mysqldb

import (
	"context"
	"fmt"
	"main/constants"
	"time"
)

type Video struct {
	Id              int64  `json:"video_id"`
	Author          int64  `json:"author_id"`
	PlayUrl         string `json:"play_url"`
	CoverUrl        string `json:"cover_url"`
	Title           string `json:"title"`
	CreateTimestamp string `json:"create_timestamp"`
}

func (u *Video) TableName() string {
	return constants.MySQLVideoTableName
}

// CreateVideo create video info
func CreateVideo(ctx context.Context, video *Video) (*Video, error) {
	video.CreateTimestamp = time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(video.Author)
	if err := MysqlDB.WithContext(ctx).Create(video).Error; err != nil {
		return video, err
	}
	return video, nil
}

// GetVideo get video info
func GetVideo(ctx context.Context, videoID int64) (*Video, error) {
	var res *Video
	if err := MysqlDB.WithContext(ctx).Where("id = ?", videoID).First(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// ListVideo list all videos info
func ListVideo(ctx context.Context) ([]*Video, error) {
	var res []*Video
	if err := MysqlDB.WithContext(ctx).Order("create_timestamp desc").Limit(constants.FeedLimit).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func ListVideoByLimit(ctx context.Context, latestTime string) ([]*Video, error) {
	var res []*Video
	if err := MysqlDB.WithContext(ctx).Where("create_timestamp <= ?", latestTime).Order("create_timestamp desc").Limit(constants.FeedLimit).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func ListVideoByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	var res []*Video

	if err := MysqlDB.WithContext(ctx).Where("author = ?", userId).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// UpdateVideo update video info
// not need
func UpdateVideoUrl(ctx context.Context, videoID int64, playUrl string, coverUrl string) error {
	err := MysqlDB.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).Update("play_url", playUrl)
	if err != nil {
		return err.Error
	}

	err = MysqlDB.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).Update("cover_url", coverUrl)
	if err != nil {
		return err.Error
	}

	return nil
}

// DeleteVideo delete video info
// not need
func DeleteVideo(ctx context.Context, videoID int64) error {
	return nil
}
