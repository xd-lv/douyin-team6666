package mysqldb

import (
	"context"
	"main/constants"
	"time"
)

type Video struct {
	Id              int64     `json:"user_id"`
	Author          int64     `json:"author_id"`
	PlayUrl         string    `json:"play_url"`
	CoverUrl        string    `json:"cover_url"`
	Title           string    `json:"title"`
	CreateTimestamp time.Time `json:"create_timestamp"`
}

func (u *Video) TableName() string {
	return constants.VideoTableName
}

// CreateVideo create video info
func CreateVideo(ctx context.Context, video *Video) error {
	video.CreateTimestamp = time.Now()
	if err := DB.WithContext(ctx).Create(video).Error; err != nil {
		return err
	}
	return nil
}

// GetVideo get video info
func GetVideo(ctx context.Context, videoID int64) (*Video, error) {
	var res *Video
	if err := DB.WithContext(ctx).Where("id = ?", videoID).First(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// ListVideo list all videos info
func ListVideo(ctx context.Context) ([]*Video, error) {
	var res []*Video

	if err := DB.WithContext(ctx).Order("create_timestamp desc").Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// UpdateVideo update video info
// not need
func UpdateVideo(ctx context.Context, videoID int64) error {
	return nil
}

// DeleteVideo delete video info
// not need
func DeleteVideo(ctx context.Context, videoID int64) error {
	return nil
}
