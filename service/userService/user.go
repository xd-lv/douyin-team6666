package userService

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	md5 "github.com/anaskhan96/go-password-encoder"
	"gorm.io/gorm"
	"main/controller"
	"main/dal/mysqldb"
	"main/utils/jwtUtil"
	"main/utils/snowflakeUtil"
	"strings"
)

func RegisterService(username string, password string) (*string, *int64, error) {
	if username == "" || password == "" {
		return nil, nil, errors.New("请检查用户名或密码是否为空")
	}
	// 查用户名是否已存在
	_, err := GetUserByUserNameService(context.TODO(), username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, err
	} else if err == nil {
		return nil, nil, errors.New("该用户名已存在")
	}
	// 密码加密
	options := &md5.Options{16, 100, 32, sha512.New}
	salt, encodePwd := md5.Encode(password, options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodePwd)
	// 注册用户
	user := &mysqldb.User{
		Id:           snowflakeUtil.NewId(),
		UserName:     username,
		UserPassword: newPassword,
	}
	if err = CreateUserService(context.TODO(), user); err != nil {
		return nil, nil, err
	}
	// 此处生成token, userId将保存在jwt中，具体实现在PayloadFunc
	token, err := CreateTokenService(user.Id)
	if err != nil {
		return nil, nil, err
	}
	return token, &user.Id, nil
}

func LoginService(username string, password string) (*string, *int64, error) {
	if username == "" || password == "" {
		return nil, nil, errors.New("请检查用户名或密码是否为空")
	}
	// 查用户
	user, err := GetUserByUserNameService(context.TODO(), username)
	if err != nil {
		return nil, nil, err
	}
	// 校验密码,密码使用md5加密
	options := &md5.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(user.UserPassword, "$")
	if check := md5.Verify(password, passwordInfo[2], passwordInfo[3], options); !check {
		return nil, nil, errors.New("密码错误")
	}
	token, err := CreateTokenService(user.Id)
	if err != nil {
		return nil, nil, err
	}
	return token, &user.Id, nil
}

func GetUserService(userId int64) (*controller.User, error) {
	user, err := mysqldb.GetUser(context.TODO(), userId)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}
	// 此处待完善社交部分
	newUser := &controller.User{
		Id:            user.Id,
		Name:          user.UserName,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	return newUser, nil
}

func CreateUserService(ctx context.Context, user *mysqldb.User) error {
	if err := mysqldb.CreateUser(ctx, user); err != nil {
		return errors.New("创建用户失败")
	}
	return nil
}

func GetUserByUserNameService(ctx context.Context, userName string) (*mysqldb.User, error) {
	user, err := mysqldb.GetUserByUserName(ctx, userName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, errors.New("通过用户名查找用户失败")
	}
	return user, nil
}

func CreateTokenService(userId int64) (*string, error) {
	token, _, err := jwtUtil.AuthMiddleware.TokenGenerator(userId)
	if err != nil {
		return nil, errors.New("生成token失败")
	}
	return &token, nil
}
