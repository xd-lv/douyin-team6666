package pack

import (
	"context"
	"main/dal/mysqldb"
)

func WithUser(userId int64) *User {
	return &User{
		Id: userId,
	}
}

func (u *User) GetUser(ctx context.Context) error {
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

func (u *User) getIsFollow(ctx context.Context) error {
	// TODO 判断是否关注该用户，需要jwt部分，可能需要添加参数
	return nil
}

func (u *User) getRocksdbUser(ctx context.Context) error {
	// TODO rocksdb 获得分数，关注数量
	return nil
}

func (u *User) getMysqlUser(ctx context.Context) error {
	dbUser, err := mysqldb.GetUser(ctx, u.Id)
	if err != nil {
		return err
	}

	u.Name = dbUser.UserName

	return nil
}
