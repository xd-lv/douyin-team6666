package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/pack"
	"main/service"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]pack.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	pack.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	pack.Response
	User pack.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token, userId, err := service.UserService.RegisterService(username, password)
	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: pack.Response{StatusCode: 0},
		UserId:   *userId,
		Token:    *token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token, userId, err := service.UserService.LoginService(username, password)
	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: pack.Response{StatusCode: 0},
		UserId:   *userId,
		Token:    *token,
	})
}

func UserInfo(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	ctx := context.Background()

	user, err := service.UserService.GetUserBody(ctx, userId)
	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: pack.Response{StatusCode: 0, StatusMsg: "success"},
		User:     *user,
	})
	return
}
