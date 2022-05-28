package mysqldb

import (
	"context"
	"testing"
)

func TestMain(m *testing.M) {
	Init()
	m.Run()
}

func TestUserCreate(t *testing.T) {
	testVideo := &Video{
		Author:   1,
		PlayUrl:  "testPlayUrl",
		CoverUrl: "testCoverUrl",
		Title:    "testTitle",
	}
	ctx := context.Background()
	testVideo, _ = CreateVideo(ctx, testVideo)
	println(testVideo.Id)
}

//
//func TestGetVideo(t *testing.T) {
//	ctx := context.Background()
//	r, err := GetVideo(ctx, 1)
//	if err != nil {
//		println(err)
//	}
//	println(r.Author)
//}
//
//func TestListVideo(t *testing.T) {
//	ctx := context.Background()
//	r, err := ListVideo(ctx)
//	if err != nil {
//		println(err)
//	}
//	println(len(r))
//}
