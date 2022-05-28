package controller

import (
	"github.com/gin-gonic/gin"
	"main/pack"
	"main/service/relationService"
	"net/http"
)

type UserListResponse struct {
	Response
	UserList []pack.User `json:"user_list"`
}

type ActionRequest struct {
	Uid   int64 `form:"user_id" json:"user_id" binding:"required"`
	ToUid int64 `form:"to_user_id" json:"to_user_id" binding:"required"`
	Type  int32 `form:"action_type" json:"action_type" binding:"required"`
}

const (
	ActionFollow       = 1
	ActionCancelFollow = 2
)

// RelationAction 关注与取消关注
func RelationAction(c *gin.Context) {
	action := ActionRequest{}
	err := c.BindQuery(&action)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "Input parameter error"})
		return
	}

	if action.Type == ActionFollow {
		err := relationService.Follow(c, action.Uid, action.ToUid)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 3, StatusMsg: "redis server error"})
			return
		}
	} else if action.Type == ActionCancelFollow {
		err := relationService.CancelFollow(c, action.Uid, action.ToUid)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 3, StatusMsg: "redis server error"})
			return
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: "parameter 'action_type' error"})
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []pack.User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []pack.User{DemoUser},
	})
}
