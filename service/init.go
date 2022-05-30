package service

import (
	"main/service/relationService"
	"main/service/userService"
	"main/service/videoService"
)

var RelationService relationService.IRelationService
var VideoService videoService.IVideoService
var UserService userService.IUserService

func init() {
	RelationService = relationService.NewRelationService()
	VideoService = videoService.NewVideoService()
	UserService = userService.NewUserService()
}
