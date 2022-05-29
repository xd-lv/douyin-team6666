package videoService

import (
	"context"
	"main/dal/miniodb"
	"main/dal/mysqldb"
	"main/pack"
	"mime/multipart"
	"strconv"
	"time"
)

func Feed(ctx context.Context, latestTime string) ([]pack.Video, int64, error) {
	var res []pack.Video
	var videos []*mysqldb.Video
	var err error
	if latestTime == "" {
		videos, err = mysqldb.ListVideo(ctx)
		if err != nil {
			return nil, 0, err
		}
	} else {

		videos, err = mysqldb.ListVideoByLimit(ctx, latestTime)
		if err != nil {
			return nil, 0, err
		}
	}

	for _, video := range videos {
		v := pack.WithVideo(video.Id)
		v.GetVideo(ctx)
		res = append(res, v)
	}
	return res, 0, nil
}

func PublishList(ctx context.Context, userId int64) ([]pack.Video, error) {
	var res []pack.Video

	videoRecordList, err := mysqldb.ListVideoByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, v := range videoRecordList {
		video := pack.WithVideo(v.Id)
		video.GetVideo(ctx)
		res = append(res, video)
	}

	return res, nil
}

func Publish(ctx context.Context, file *multipart.FileHeader, title string, userId int64) error {
	user, err := mysqldb.GetUser(ctx, userId)
	if err != nil {
		return err
	}

	videoRecord := &mysqldb.Video{
		Author:          userId,
		PlayUrl:         "",
		CoverUrl:        "",
		Title:           title,
		CreateTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	videoRecord, err = mysqldb.CreateVideo(ctx, videoRecord)
	if err != nil {
		return err
	}

	videos, err := mysqldb.ListVideoByUserId(ctx, userId)
	if err != nil {
		return err
	}

	var purl, curl string
	if len(videos) == 0 {
		purl, curl, err = upload(file, user.UserName, videoRecord.Id, false)
	} else {
		purl, curl, err = upload(file, user.UserName, videoRecord.Id, true)
	}
	if err != nil {
		return err
	}

	err = mysqldb.UpdateVideoUrl(ctx, videoRecord.Id, purl, curl)
	if err != nil {
		return err
	}

	return nil
}

func upload(sfile *multipart.FileHeader, userName string, videoId int64, isCreateBucket bool) (string, string, error) {
	var purl, curl string
	var bucketName = userName + "bucket"
	file, err := sfile.Open()
	if err != nil {
		return purl, curl, err
	}

	if isCreateBucket {
		// 需要创建bucket
		err := miniodb.Create(bucketName)
		if err != nil {
			return purl, curl, err
		}
	}
	purl, err = miniodb.Upload(file, bucketName, strconv.FormatInt(videoId, 10)+"-play", sfile.Size)
	if err != nil {
		return purl, curl, err
	}
	// TODO cover 截取
	return purl, curl, nil
}
