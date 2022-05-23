package main

import (
	"github.com/gin-gonic/gin"
	"main/dal/mysqldb"
)

func main() {
	// 初始化工作
	mysqldb.Init()
	r := gin.Default()

	initRouter(r)

	r.Run("127.0.0.1:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
