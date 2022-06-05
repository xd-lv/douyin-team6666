package mysqldb

import (
	"context"
	"main/constants"
)

type UserInfo struct {
	Id           int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

func (u *UserInfo) TableName() string {
	return constants.MySQLUserTableName
}

// CreateUser create user info
func CreateUser(ctx context.Context, user *UserInfo) error {
	if err := MysqlDB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUser get user info
func GetUser(ctx context.Context, userID int64) (*UserInfo, error) {
	var res *UserInfo
	if err := MysqlDB.WithContext(ctx).Where("id = ?", userID).First(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// ListUser list all users info
func ListUser(ctx context.Context) ([]*UserInfo, error) {
	var res []*UserInfo

	if err := MysqlDB.WithContext(ctx).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// UpdateUser update user info
// not need
func UpdateUser(ctx context.Context, userID int64, userNewName, userNewPassword *string) error {
	return nil
}

// DeleteUser delete user info
// not need
func DeleteUser(ctx context.Context, userID int64) error {
	return nil
}

// GetUserByUserName Get user info by username
func GetUserByUserName(ctx context.Context, userName string) (*UserInfo, error) {
	var res *UserInfo
	if err := MysqlDB.WithContext(ctx).Where("user_name = ?", userName).First(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// GetUserNameByID 通过 userID 获取对应的用户名
func GetUserNameByID(ctx context.Context, userID int64) (string, error) {
	var user UserInfo
	if err := MysqlDB.WithContext(ctx).Select("user_name").First(&user, userID).Error; err != nil {
		return "", constants.ErrMysqlServer
	}
	return user.UserName, nil
}
