package favoriteService

import (
	"context"
	"testing"
)

func TestImpl_IsFavorite(t *testing.T) {
	service := NewIFavoriteService()
	background := context.Background()
	ret, err := service.IsFavorite(background, 536147966373662722, 4)
	if err != nil {
		return
	}

	println(ret)
}
