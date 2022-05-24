package dal

import (
	"main/dal/miniodb"
	"main/dal/mysqldb"
)

// Init init dal
func Init() {
	mysqldb.Init() // mysql

	// TODO init() rocksdb

	// TODO init() minio
	miniodb.Init()
}
