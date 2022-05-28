package dal

import (
	"main/dal/miniodb"
	"main/dal/mysqldb"
)

// Init init dal
func Init() {
	mysqldb.Init() // mysql
	miniodb.Init() // minio
}
