package service

import (
	"main/service/relationService"
)

var RelationService relationService.IRelationService

func init() {
	RelationService = relationService.NewRelationService()
}
