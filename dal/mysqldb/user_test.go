package mysqldb

import (
	"context"
	"testing"
)

func TestCreateUser(t *testing.T) {
	testUser := &User{
		UserName:     "test1",
		UserPassword: "test",
	}
	ctx := context.Background()
	CreateUser(ctx, testUser)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	u, err := GetUser(ctx, 1)
	if err != nil {
		println(err)
	}
	println(u.UserName)
}

func TestListUser(t *testing.T) {
	ctx := context.Background()
	list, err := ListUser(ctx)
	if err != nil {
		println(err)
	}
	println(len(list))
}
