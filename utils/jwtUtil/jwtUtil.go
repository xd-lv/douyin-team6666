package jwtUtil

import (
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"main/constants"
	"net/http"
	"strconv"
	"time"
)

var (
	AuthMiddleware *jwt.GinJWTMiddleware // JWT中间件
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func init() {
	AuthMiddleware, _ = NewJwt()
}

func NewJwt() (*jwt.GinJWTMiddleware, error) {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:         []byte(constants.JWTKey),
		Timeout:     time.Hour * 24 * 7,
		MaxRefresh:  time.Hour,
		IdentityKey: constants.IdentityKey,

		Authenticator: func(c *gin.Context) (interface{}, error) { // 根据登录信息对用户进行身份验证的回调函数
			username := c.Query("username")
			password := c.Query("password")
			c.Set("userId", 123456)
			if username != "admin" || password != "123456" {
				return nil, jwt.ErrFailedAuthentication
			}
			return username, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) { // 处理不进行授权的逻辑
			if code != http.StatusOK {
				code = 1
			} else {
				code = 0
			}
			c.JSON(http.StatusOK, Response{
				StatusCode: int32(code),
				StatusMsg:  message,
			})
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims { // 登录期间的回调函数
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{ // int64赋值给interface类型后转成float64，精度丢失
					constants.IdentityKey: strconv.FormatInt(v, 10),
				}
			}
			return jwt.MapClaims{}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			userId, exist := c.Get("userId")
			if !exist {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					StatusMsg:  "无法获取到用户id",
				})
				return
			}
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   int64(userId.(int)),
				Token:    token,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return authMiddleware, err
}

// AuthTokenForm 表单提交token
func AuthTokenForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		c.Request.Header.Set("Authorization", "Bearer "+token)
		fmt.Println(token)
		c.Next()
	}
}

// GetLoginUserId 获取当前登录用户的userId
func GetLoginUserId(c *gin.Context) (userId int64, err error) {
	claim, err := AuthMiddleware.GetClaimsFromJWT(c)
	if err != nil {
		return 0, errors.New("Failed to get user id")
	}
	userIdStr := claim[constants.IdentityKey].(string)
	userId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, errors.New("Failed to convert string to int64")
	}
	return userId, nil
}
