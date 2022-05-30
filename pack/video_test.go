package pack

import (
	"context"
	"testing"
)

func TestGetVideo(t *testing.T) {
	ctx := context.Background()
	video := WithVideo(1)
	video.GetVideo(ctx)
	println(video.Author.Id)
}
