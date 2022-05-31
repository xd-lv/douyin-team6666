package commentService

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"main/dal/mysqldb"
	"main/dal/redisdb"
	"main/pack"
	//"main/service"

	"main/utils/snowflakeUtil"
	"strconv"
	"strings"
	"time"
)

var rdbComment = redisdb.RDBComment
var rdbRelation = redisdb.RDB

type ICommentService interface {
	CreateComment(ctx context.Context, videoId int64, userId int64, commentText string) error
	DeleteComment(ctx context.Context, videoId int64, commentId int64) error
	ListComment(ctx context.Context, videoId int64) ([]pack.Comment, error)
}

type Impl struct {
}

func NewCommentService() ICommentService {
	return &Impl{}
}

// Comment 评论插入操作
func (comment *Impl) CreateComment(ctx context.Context, videoId int64, userId int64, commentText string) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	commentIdStr := strconv.FormatInt(snowflakeUtil.NewId(), 10)
	userIDStr := strconv.FormatInt(userId, 10)

	commentIdAndText := fmt.Sprintf("%s+%s+%s+%s", commentIdStr, timestamp, userIDStr, commentText)

	fmt.Println(commentIdAndText)
	videoIdStr := strconv.FormatInt(videoId, 10)

	pipe := rdbComment.TxPipeline()

	pipe.LPush(ctx, videoIdStr, commentIdAndText)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errors.New("redis server error. Method: CreateComment")
	}
	return nil
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
		fmt.Println(strSplit[0])
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

		tempUser, err := comment.GetRelationAuthor(ctx, video.Author, tempUserId)

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

// GetRelationAuthor 获取 `作者`(视频作者、评论作者等) 的关系信息，其中还包含 `用户` 对 `作者` 的关注情况
func (comment *Impl) GetRelationAuthor(ctx context.Context, authorId int64, userId int64) (pack.User, error) {
	user, err := comment.GetRelationUser(ctx, authorId)
	if err != nil {
		return user, err
	}

	user.IsFollow, err = comment.isAFollowB(ctx, userId, authorId)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetRelationUser 获取 `用户`(登录用户|操作人) 的关系信息，包含关注数和粉丝数，字段 `IsFollow` 默认 false
func (comment *Impl) GetRelationUser(ctx context.Context, userId int64) (pack.User, error) {
	user := pack.User{Id: userId}

	username, err := mysqldb.GetUserNameByID(ctx, user.Id)
	if err != nil {
		return user, err
	}
	user.Name = username

	user.FollowCount, err = comment.getFollowCount(ctx, user.Id)
	if err != nil {
		return user, err
	}
	user.FollowerCount, err = comment.getFollowerCount(ctx, user.Id)
	if err != nil {
		return user, err
	}

	return user, nil
}

// getFollowCount 获取关注数量
func (comment *Impl) getFollowCount(ctx context.Context, userId int64) (int64, error) {
	followKey := "follow_" + strconv.FormatInt(userId, 10)
	followCount, err := rdbRelation.ZCard(ctx, followKey).Result()
	if err != nil {
		return 0, errors.New("redis server error")
	}
	return followCount, nil
}

// getFollowerCount 获取粉丝数量
func (comment *Impl) getFollowerCount(ctx context.Context, userId int64) (int64, error) {
	fansKey := "fans_" + strconv.FormatInt(userId, 10)
	followerCount, err := rdbRelation.ZCard(ctx, fansKey).Result()
	if err != nil {
		return 0, errors.New("redis server error")
	}
	return followerCount, nil
}

// isAFollowB 判断用户A是否关注了用户B
func (comment *Impl) isAFollowB(ctx context.Context, userAId int64, userBId int64) (bool, error) {
	followKey := "follow_" + strconv.FormatInt(userAId, 10)
	userBIdStr := strconv.FormatInt(userBId, 10)
	if err := rdbRelation.ZRank(ctx, followKey, userBIdStr).Err(); err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, errors.New("redis server error")
	}
	return true, nil
}
