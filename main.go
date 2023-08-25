package main

import (
	"log"

	"github.com/CheemsGoUp/Simplified-Douyin-Project/controller"
	"github.com/CheemsGoUp/Simplified-Douyin-Project/global"
	"github.com/CheemsGoUp/Simplified-Douyin-Project/service"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	go service.RunMessageServer()

	var err error
	global.DB, err = controller.InitDB()
	if err != nil {
		log.Fatalf("fail to initiate database, %v\n", err)
	}

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
