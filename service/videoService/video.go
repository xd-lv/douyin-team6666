package videoService

import (
	"context"
	"main/dal/miniodb"
	"main/dal/mysqldb"
	"main/pack"
	"main/service/favoriteService"
	"main/service/userService"
	"mime/multipart"
	"strconv"
	"time"
)

type IVideoService interface {
	Feed(ctx context.Context, latestTime string) ([]pack.Video, int64, error)
	PublishList(ctx context.Context, userId int64) ([]pack.Video, error)
	Publish(ctx context.Context, file *multipart.FileHeader, title string, userId int64) error
	GetVideoBody(ctx context.Context, videoId int64) (*pack.Video, error)
}

type Impl struct {
	userService     userService.IUserService
	favoriteService favoriteService.IFavoriteService
}

func NewVideoService() IVideoService {
	return &Impl{
		userService:     userService.NewUserService(),
		favoriteService: favoriteService.NewIFavoriteService(),
	}
}

func (vs *Impl) Feed(ctx context.Context, latestTime string) ([]pack.Video, int64, error) {
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
		v, err := vs.GetVideoBody(ctx, video.Id)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, *v)
	}
	return res, 0, nil
}

func (vs *Impl) PublishList(ctx context.Context, userId int64) ([]pack.Video, error) {
	var res []pack.Video

	videoRecordList, err := mysqldb.ListVideoByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, v := range videoRecordList {
		video, err := vs.GetVideoBody(ctx, v.Id)
		if err != nil {
			return nil, err
		}
		res = append(res, *video)
	}

	return res, nil
}

func (vs *Impl) Publish(ctx context.Context, file *multipart.FileHeader, title string, userId int64) error {
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

	var purl, curl string

	purl, curl, err = upload(file, user.UserName, videoRecord.Id)

	if err != nil {
		return err
	}

	err = mysqldb.UpdateVideoUrl(ctx, videoRecord.Id, purl, curl)
	if err != nil {
		return err
	}

	return nil
}

func (vs *Impl) GetVideoBody(ctx context.Context, videoId int64) (*pack.Video, error) {
	res := &pack.Video{
		Id: videoId,
	}

	authorId, err := vs.getMysqlVideo(ctx, res)
	if err != nil {
		return res, err
	}

	author, err := vs.userService.GetUserBody(ctx, authorId)
	if err != nil {
		return res, err
	}
	res.Author = *author

	// TODO comment&favorite
	fCount, err := vs.favoriteService.GetFavoriteCount(ctx, videoId)
	res.FavoriteCount = fCount

	return res, nil
}

func (vs *Impl) getMysqlVideo(ctx context.Context, videoBody *pack.Video) (int64, error) {

	dbVideo, err := mysqldb.GetVideo(ctx, videoBody.Id)
	if err != nil {
		return dbVideo.Author, err
	}

	videoBody.Title = dbVideo.Title
	videoBody.CoverUrl = dbVideo.CoverUrl
	videoBody.PlayUrl = dbVideo.PlayUrl

	author, err := vs.userService.GetUserBody(ctx, dbVideo.Author)
	if err != nil {
		return dbVideo.Author, err
	}

	videoBody.Author = *author

	return dbVideo.Author, nil
}

func upload(sfile *multipart.FileHeader, userName string, videoId int64) (string, string, error) {
	var purl, curl string
	var bucketName = userName + "bucket"
	file, err := sfile.Open()
	if err != nil {
		return purl, curl, err
	}

	exists := miniodb.IsExists(bucketName)

	if !exists {
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
