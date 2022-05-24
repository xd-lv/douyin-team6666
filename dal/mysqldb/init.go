package mysqldb

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
	"main/constants"
)

var MysqlDB *gorm.DB

// Init init DB
func Init() {
	var err error

	MysqlDB, err = gorm.Open(mysql.Open(constants.MySQLDefaultDSN),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}

	if err = MysqlDB.Use(gormopentracing.New()); err != nil {
		panic(err)
	}

}
