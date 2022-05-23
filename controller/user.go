package controller

import (
	"context"
	"crypto/sha512"
	"fmt"
	md5 "github.com/anaskhan96/go-password-encoder"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/dal/mysqldb"
	"main/utils/jwtUtil"
	"main/utils/snowflakeUtil"
	"net/http"
	"strings"
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
	if username == "" || password == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "请输入用户名和密码",
		})
		return
	}
	// 查用户名是否已存在
	_, err := mysqldb.GetUserByUserName(context.TODO(), username)
	if err == nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户名已存在",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "注册失败",
		})
		return
	}
	// 密码加密
	options := &md5.Options{16, 100, 32, sha512.New}
	salt, encodePwd := md5.Encode(password, options)
	password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodePwd)
	// 注册用户
	user := &mysqldb.User{
		Id:           snowflakeUtil.NewId(),
		UserName:     username,
		UserPassword: password,
	}
	if err = mysqldb.CreateUser(context.TODO(), user); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "注册失败",
		})
		return
	}
	// 此处生成token, userId将保存在jwt中，具体实现在PayloadFunc
	token, _, err := jwtUtil.AuthMiddleware.TokenGenerator(user.Id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	})
	return
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "请输入用户名和密码",
		})
		return
	}
	// 通过用户名查找
	user, err := mysqldb.GetUserByUserName(context.TODO(), username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "没有该用户",
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 校验密码,密码使用md5加密
	options := &md5.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(user.UserPassword, "$")
	if check := md5.Verify(password, passwordInfo[2], passwordInfo[3], options); !check {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "密码错误",
		})
		return
	}
	// 此处生成token, userId将保存在jwt中，具体实现在PayloadFunc
	token, _, err := jwtUtil.AuthMiddleware.TokenGenerator(user.Id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
