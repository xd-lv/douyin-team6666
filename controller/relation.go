package controller

import (
	"github.com/gin-gonic/gin"
	"main/pack"
	"main/service"
	"net/http"
	"strconv"
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
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Input parameter error"})
		return
	}

	if action.Type == ActionFollow {
		err := service.RelationService.Follow(c, action.Uid, action.ToUid)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: err.Error()})
			return
		}
	} else if action.Type == ActionCancelFollow {
		err := service.RelationService.CancelFollow(c, action.Uid, action.ToUid)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 2, StatusMsg: err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "parameter 'action_type' error"})
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// FollowList 关注列表
func FollowList(c *gin.Context) {
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Input parameter error",
			},
			UserList: []pack.User{},
		})
		return
	}

	users, err := service.RelationService.GetFollowList(c, uid)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 2,
				StatusMsg:  err.Error(),
			},
			UserList: users,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Input parameter error",
			},
			UserList: []pack.User{},
		})
		return
	}

	users, err := service.RelationService.GetFollowerList(c, uid)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 2,
				StatusMsg:  err.Error(),
			},
			UserList: users,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}
