package dal

import "main/dal/mysqldb"

// Init init dal
func Init() {
	mysqldb.Init() // mysql

	// TODO init() rocksdb

	// TODO init() minio
}
