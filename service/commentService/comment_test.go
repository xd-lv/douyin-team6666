package commentService

import (
	"context"
	"fmt"
	"main/dal/mysqldb"
	"testing"
)

func TestImpl_CreateComment(t *testing.T) {
	service := NewCommentService()
	ctx := context.Background()
	str := "饿了饿啊啊啊啊啊啊啊啊"
	err := service.CreateComment(ctx, 1, 536614587361923072, str)
	if err != nil {
		println(err.Error())
		return
	}
	println("ok")
}

func TestImpl_ListComment(t *testing.T) {
	mysqldb.Init()
	service := NewCommentService()
	ctx := context.Background()
	comment, err := service.ListComment(ctx, 1)
	if err != nil {
		return
	}
	for _, p := range comment {
		fmt.Println(p.Id, p.User, p.Content, p.CreateDate)
	}
	println("ok")
}

func TestImpl_DeleteComment(t *testing.T) {
	mysqldb.Init()
	service := NewCommentService()
	background := context.Background()
	err := service.DeleteComment(background, 1, 537765995356360704)
	if err != nil {
		return
	}
	comment, err := service.ListComment(background, 1)
	if err != nil {
		return
	}
	for _, p := range comment {
		fmt.Println(p.Id, p.User, p.Content, p.CreateDate)
	}
	println("ok")
}

/*
userid             name

536147966373662720 admin
536614587361923072 xxd@qq.com
536883152367390720 Gang
537321358355337216 lvxiaodong
537325256709246976 lxd

video id

7
6
5
4
3
2
1
*/
