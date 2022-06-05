package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/constants"
	"main/pack"
	"main/service"
	"net/http"
	"strconv"
)

type likeListResponse struct {
	pack.Response
	videoList []pack.Video `json:"video_list"`
}

type favoriteActionRequest struct {
	Uid  int64 `form:"user_id" json:"user_id" binding:"required"`
	Vid  int64 `form:"video_id" json:"video_id" binding:"required"`
	Type int32 `form:"action_type" json:"action_type" binding:"required"`
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	action := favoriteActionRequest{}
	err := c.BindQuery(&action)

	if err != nil {
		c.JSON(http.StatusOK, pack.Response{StatusCode: 1, StatusMsg: constants.ErrInvalidParams.Error()})
		return
	}
	if action.Type == constants.ActionLike {
		err := service.FavoriteService.Like(c, action.Uid, action.Vid)
		if err != nil {
			c.JSON(http.StatusOK, pack.Response{StatusCode: 2, StatusMsg: err.Error()})
			return
		}
	} else if action.Type == constants.ActionUnlike {
		err := service.FavoriteService.Unlike(c, action.Uid, action.Vid)
		if err != nil {
			c.JSON(http.StatusOK, pack.Response{StatusCode: 2, StatusMsg: err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, pack.Response{StatusCode: 1, StatusMsg: constants.ErrInvalidParams.Error()})
	}

	c.JSON(http.StatusOK, pack.Response{StatusCode: 0})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	videos := make([]pack.Video, 0)
	if err != nil {
		c.JSON(http.StatusOK, likeListResponse{
			Response: pack.Response{StatusCode: 1, StatusMsg: constants.ErrInvalidParams.Error()},
		},
		)
	}

	vIds, err := service.FavoriteService.GetLikeList(c, uid)
	if err != nil {
		c.JSON(http.StatusOK, likeListResponse{
			Response: pack.Response{StatusCode: 1, StatusMsg: err.Error()},
		},
		)
	}

	for _, id := range vIds {
		temp, err := service.VideoService.GetVideoBody(c, id)
		if err != nil {
			c.JSON(http.StatusOK, likeListResponse{
				Response: pack.Response{StatusCode: 1, StatusMsg: err.Error()},
			},
			)
		}

		video := pack.Video{Id: temp.Id,
			Title:         temp.Title,
			Author:        temp.Author,
			PlayUrl:       temp.PlayUrl,
			CoverUrl:      temp.CoverUrl,
			FavoriteCount: temp.FavoriteCount,
			CommentCount:  0,
			IsFavorite:    true}
		fmt.Println(video)
		videos = append(videos, video)
	}

	c.JSON(http.StatusOK, likeListResponse{
		Response: pack.Response{
			StatusCode: 0,
		},
		videoList: videos,
	})
}
