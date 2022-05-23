package jwtUtil

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"main/constants"
	"main/controller"
	"net/http"
	"time"
)

var (
	AuthMiddleware *jwt.GinJWTMiddleware // JWT中间件
)

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
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: int32(code),
				StatusMsg:  message,
			})
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims { // 登录期间的回调函数
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					constants.IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			userId, exist := c.Get("userId")
			if !exist {
				c.JSON(http.StatusOK, controller.Response{
					StatusCode: 1,
					StatusMsg:  "无法获取到用户id",
				})
				return
			}
			c.JSON(http.StatusOK, controller.UserLoginResponse{
				Response: controller.Response{StatusCode: 0},
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
