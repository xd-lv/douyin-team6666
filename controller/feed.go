package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/pack"
	"main/service"
	"net/http"
)

type FeedResponse struct {
	pack.Response
	VideoList []pack.Video `json:"video_list,omitempty"`
	NextTime  int64        `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	ctx := context.Background()
	latestTime := c.Query("latest_time")
	videoList, nextTime, err := service.VideoService.Feed(ctx, latestTime)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: pack.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  pack.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}
