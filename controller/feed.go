package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/pack"
	"main/service/videoService"
	"net/http"
)

type FeedResponse struct {
	Response
	VideoList []pack.Video `json:"video_list,omitempty"`
	NextTime  int64        `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	ctx := context.Background()
	latestTime := c.Query("latest_time")
	videoList, nextTime, err := videoService.Feed(ctx, latestTime)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}
