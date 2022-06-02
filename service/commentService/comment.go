package commentService

import (
	"context"
	"errors"
	"fmt"
	"main/dal/mysqldb"
	"main/dal/redisdb"
	"main/pack"
	"main/service/relationService"

	//"main/service"

	"main/utils/snowflakeUtil"
	"strconv"
	"strings"
	"time"
)

var rdbComment = redisdb.RDBComment

type ICommentService interface {
	CreateComment(ctx context.Context, videoId int64, userId int64, commentText string) (pack.Comment, error)
	DeleteComment(ctx context.Context, videoId int64, commentId int64) error
	ListComment(ctx context.Context, videoId int64) ([]pack.Comment, error)
}

type Impl struct {
	relationService relationService.IRelationService
}

func NewCommentService() ICommentService {
	return &Impl{
		relationService: relationService.NewRelationService(),
	}
}

// Comment 评论插入操作
func (comment *Impl) CreateComment(ctx context.Context, videoId int64, userId int64, commentText string) (pack.Comment, error) {
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	commentId := snowflakeUtil.NewId()
	commentIdStr := strconv.FormatInt(commentId, 10)
	userIDStr := strconv.FormatInt(userId, 10)

	commentIdAndText := fmt.Sprintf("%s+%s+%s+%s", commentIdStr, timestampStr, userIDStr, commentText)

	videoIdStr := strconv.FormatInt(videoId, 10)

	pipe := rdbComment.TxPipeline()

	pipe.LPush(ctx, videoIdStr, commentIdAndText)

	_, err := pipe.Exec(ctx)
	if err != nil {
		fmt.Println("redis server error. Method: CreateComment, exec")
		return pack.Comment{}, errors.New("redis server error. Method: CreateComment, exec")
	}

	video, err := mysqldb.GetVideo(ctx, videoId)

	if err != nil {
		fmt.Println("redis server error. Method: CreateComment, video")
		return pack.Comment{}, errors.New("redis server error. Method: CreateComment, video")
	}

	user, err := comment.relationService.GetRelationAuthor(ctx, video.Author, userId)

	if err != nil {
		return pack.Comment{}, errors.New("redis server error. Method: CreateComment, user")
	}

	timeStampStrData := time.Unix(timestamp, 0).Format("01/02")

	return pack.Comment{
		Id:         commentId,
		User:       user,
		Content:    commentText,
		CreateDate: timeStampStrData,
	}, nil
}

// 评论删除操作
func (comment *Impl) DeleteComment(ctx context.Context, videoId int64, commentId int64) error {
	videoIdStr := strconv.FormatInt(videoId, 10)
	commentIdStr := strconv.FormatInt(commentId, 10)

	resultComment, err := rdbComment.LRange(ctx, videoIdStr, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, tempCommentStr := range resultComment {
		strSplit := strings.SplitN(tempCommentStr, "+", 4)
		if strings.Compare(strSplit[0], commentIdStr) == 0 {
			commentIdStr = tempCommentStr
			break
		}
	}

	pipe := rdbComment.TxPipeline()

	pipe.LRem(ctx, videoIdStr, 0, commentIdStr)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// 评论列表

func (comment *Impl) ListComment(ctx context.Context, videoId int64) ([]pack.Comment, error) {

	comments := make([]pack.Comment, 0)
	videoIdStr := strconv.FormatInt(videoId, 10)

	lRange, err := rdbComment.LRange(ctx, videoIdStr, 0, -1).Result()
	if err != nil {
		return comments, errors.New("redis server error. Method: ListComment")
	}
	video, _ := mysqldb.GetVideo(ctx, videoId)

	for _, tempCommentStr := range lRange {
		strSplit := strings.SplitN(tempCommentStr, "+", 4)
		tempCommentId, _ := strconv.ParseInt(strSplit[0], 10, 64)
		tempTimeStamp, _ := strconv.ParseInt(strSplit[1], 10, 64)
		tempUserId, _ := strconv.ParseInt(strSplit[2], 10, 64)

		tempUser, err := comment.relationService.GetRelationAuthor(ctx, video.Author, tempUserId)

		if err != nil {
			return []pack.Comment{}, err
		}

		tempTimeStampStr := time.Unix(tempTimeStamp, 0).Format("01/02")

		tempComment := pack.Comment{
			Id:         tempCommentId,
			User:       tempUser,
			Content:    strSplit[3],
			CreateDate: tempTimeStampStr,
		}
		comments = append(comments, tempComment)
	}
	return comments, nil
}
