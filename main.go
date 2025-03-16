package main

import (
	// "net/http"

	"github.com/gin-gonic/gin"

	"os"

	"Community_Notification_System/configs"
	"Community_Notification_System/database"
	"Community_Notification_System/middlewares"
	routers "Community_Notification_System/routers"
)

func main() {
	// 初始化 Gin
	router := gin.Default()

	//env 初始化
	configs.InitConfig()

	//DB 初始化
	database.InitDB()

	//跨域的 Middleware
	router.Use(middlewares.CORSMiddleware())

	apiGroup := router.Group("/api")
	routers.RegisterRoutes(apiGroup)
	// 啟動伺服器

	port := os.Getenv("PORT")
	router.Run(port)
}
