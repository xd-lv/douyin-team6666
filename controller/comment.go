package controller

import (
	"github.com/gin-gonic/gin"
	"main/pack"
	"net/http"
)

type CommentListResponse struct {
	pack.Response
	CommentList []pack.Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, pack.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, pack.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    pack.Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
