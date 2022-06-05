package service

import (
	"main/service/commentService"
	"main/service/favoriteService"
	"main/service/relationService"
	"main/service/userService"
	"main/service/videoService"
)

var RelationService relationService.IRelationService
var VideoService videoService.IVideoService
var UserService userService.IUserService
var CommentService commentService.ICommentService
var FavoriteService favoriteService.IFavoriteService

func init() {
	RelationService = relationService.NewRelationService()
	VideoService = videoService.NewVideoService()
	UserService = userService.NewUserService()
	CommentService = commentService.NewCommentService()
	FavoriteService = favoriteService.NewIFavoriteService()
}
