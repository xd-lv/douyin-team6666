package userService

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	md5 "github.com/anaskhan96/go-password-encoder"
	"gorm.io/gorm"
	"main/dal/mysqldb"
	"main/pack"
	"main/service/relationService"
	"main/utils/jwtUtil"
	"main/utils/snowflakeUtil"
	"strings"
)

type IUserService interface {
	RegisterService(username string, password string) (*string, *int64, error)
	LoginService(username string, password string) (*string, *int64, error)
	GetUserService(ctx context.Context, userId int64) (*mysqldb.UserInfo, error)
	GetUserBody(ctx context.Context, userId int64) (*pack.User, error)
}

type Impl struct {
	relationService relationService.IRelationService
}

func NewUserService() IUserService {
	return &Impl{
		relationService: relationService.NewRelationService(),
	}
}

func (us *Impl) RegisterService(username string, password string) (*string, *int64, error) {
	if username == "" || password == "" {
		return nil, nil, errors.New("请检查用户名或密码是否为空")
	}
	// 查用户名是否已存在
	_, err := us.getUserByUserNameService(context.TODO(), username)
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
	user := &mysqldb.UserInfo{
		Id:           snowflakeUtil.NewId(),
		UserName:     username,
		UserPassword: newPassword,
	}
	if err = us.createUserService(context.TODO(), user); err != nil {
		return nil, nil, err
	}
	// 此处生成token, userId将保存在jwt中，具体实现在PayloadFunc
	token, err := us.createTokenService(user.Id)
	if err != nil {
		return nil, nil, err
	}
	return token, &user.Id, nil
}

func (us *Impl) LoginService(username string, password string) (*string, *int64, error) {
	if username == "" || password == "" {
		return nil, nil, errors.New("请检查用户名或密码是否为空")
	}
	// 查用户
	user, err := us.getUserByUserNameService(context.TODO(), username)
	if err != nil {
		return nil, nil, err
	}
	// 校验密码,密码使用md5加密
	options := &md5.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(user.UserPassword, "$")
	if check := md5.Verify(password, passwordInfo[2], passwordInfo[3], options); !check {
		return nil, nil, errors.New("密码错误")
	}
	token, err := us.createTokenService(user.Id)
	if err != nil {
		return nil, nil, err
	}
	return token, &user.Id, nil
}

func (us *Impl) GetUserService(ctx context.Context, userId int64) (*mysqldb.UserInfo, error) {
	user, err := mysqldb.GetUser(ctx, userId)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}
	return user, nil
}

func (us *Impl) GetUserBody(ctx context.Context, userId int64) (*pack.User, error) {
	res := &pack.User{
		Id: userId,
	}

	userInfo, err := us.GetUserService(ctx, userId)
	if err != nil {
		return nil, err
	}
	res.Name = userInfo.UserName

	//us.relationService

	return res, nil
}

func (us *Impl) createUserService(ctx context.Context, user *mysqldb.UserInfo) error {
	if err := mysqldb.CreateUser(ctx, user); err != nil {
		return errors.New("创建用户失败")
	}
	return nil
}

func (us *Impl) getUserByUserNameService(ctx context.Context, userName string) (*mysqldb.UserInfo, error) {
	user, err := mysqldb.GetUserByUserName(ctx, userName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, errors.New("通过用户名查找用户失败")
	}
	return user, nil
}

func (us *Impl) createTokenService(userId int64) (*string, error) {
	token, _, err := jwtUtil.AuthMiddleware.TokenGenerator(userId)
	if err != nil {
		return nil, errors.New("生成token失败")
	}
	return &token, nil
}
