package pack

import (
	"context"
	"main/dal"
	"testing"
)

func TestMain(m *testing.M) {
	dal.Init()
	m.Run()
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	user := WithUser(1)
	user.GetUser(ctx)
	println(user.followCount)
}
