package mysqldb

import (
	"context"
	"main/constants"
)

type User struct {
	Id           int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

func (u *User) TableName() string {
	return constants.MySQLUserTableName
}

// CreateUser create user info
func CreateUser(ctx context.Context, user *User) error {
	if err := MysqlDB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUser get user info
func GetUser(ctx context.Context, userID int64) (*User, error) {
	var res *User
	if err := MysqlDB.WithContext(ctx).Where("id = ?", userID).First(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// ListUser list all users info
func ListUser(ctx context.Context) ([]*User, error) {
	var res []*User

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
func GetUserByUserName(ctx context.Context, userName string) (*User, error) {
	var res *User
	if err := MysqlDB.WithContext(ctx).Where("user_name = ?", userName).First(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}
