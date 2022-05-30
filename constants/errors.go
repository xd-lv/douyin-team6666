package constants

import "errors"

var (
	ErrInvalidParams = errors.New("input parameter error")
	ErrMysqlServer   = errors.New("mysql server error")
	ErrRedisServer   = errors.New("redis server error")
)
