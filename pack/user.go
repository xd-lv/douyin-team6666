package pack

import (
	"context"
	"main/dal/mysqldb"
)

type user struct {
	id           int64
	name         string
	followCount  int64
	folloerCount int64
	isFollow     bool
}

func WithUser(userId int64) *user {
	return &user{
		id: userId,
	}
}

func (u *user) GetUser(ctx context.Context) error {
	err := u.getMysqlUser(ctx)
	if err != nil {
		return err
	}

	err = u.getRocksdbUser(ctx)
	if err != nil {
		return err
	}

	err = u.getIsFollow(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) getIsFollow(ctx context.Context) error {
	// TODO 判断是否关注该用户，需要jwt部分，可能需要添加参数
	return nil
}

func (u *user) getRocksdbUser(ctx context.Context) error {
	// TODO rocksdb 获得分数，关注数量
	return nil
}

func (u *user) getMysqlUser(ctx context.Context) error {
	dbUser, err := mysqldb.GetUser(ctx, u.id)
	if err != nil {
		return err
	}

	u.name = dbUser.UserName

	return nil
}
