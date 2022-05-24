package controller

import (
	"github.com/gin-gonic/gin"
	"main/service/userService"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token, userId, err := userService.RegisterService(username, password)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   *userId,
		Token:    *token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token, userId, err := userService.LoginService(username, password)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   *userId,
		Token:    *token,
	})
}

func UserInfo(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	user, err := userService.GetUserService(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0, StatusMsg: "success"},
		User: User{
			Id:            user.Id,
			Name:          user.UserName,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		},
	})
	return
}
