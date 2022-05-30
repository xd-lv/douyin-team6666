package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/constants"
	"main/pack"
	"main/service"
	"main/utils/jwtUtil"
	"net/http"
	"strconv"
)

type VideoListResponse struct {
	pack.Response
	VideoList []pack.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	title := c.PostForm("title")

	claim, err := jwtUtil.AuthMiddleware.GetClaimsFromJWT(c)
	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1, StatusMsg: "User doesn't exist",
		})
		return
	}
	userId, _ := strconv.ParseInt(claim[constants.IdentityKey].(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1, StatusMsg: "User doesn't exist",
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1, StatusMsg: err.Error(),
		})
		return
	}

	ctx := context.Background()

	err = service.VideoService.Publish(ctx, data, title, userId)
	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1, StatusMsg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, pack.Response{
		StatusCode: 1, StatusMsg: "uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {

	claim, err := jwtUtil.AuthMiddleware.GetClaimsFromJWT(c)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: pack.Response{
				StatusCode: 1, StatusMsg: "User doesn't exist",
			},
		})
		return
	}
	userId, _ := strconv.ParseInt(claim[constants.IdentityKey].(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: pack.Response{
				StatusCode: 1, StatusMsg: "User doesn't exist",
			},
		})
		return
	}

	ctx := context.Background()

	videoList, err := service.VideoService.PublishList(ctx, userId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: pack.Response{
				StatusCode: 1, StatusMsg: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: pack.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
