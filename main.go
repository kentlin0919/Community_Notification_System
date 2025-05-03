package main

import (
	// "net/http"

	"github.com/gin-gonic/gin"

	"os"

	"Community_Notification_System/configs"
	"Community_Notification_System/database"
	"Community_Notification_System/middlewares"
	routers "Community_Notification_System/routers"

	// swaggerFiles 用於提供 Swagger UI 所需的檔案
	swaggerFiles "github.com/swaggo/files"
	// ginSwagger 用於在 Gin 中嵌入 Swagger UI
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Community_Notification_System/docs" // swagger 產生的 docs package
)

// @title           Community_Notification_System
// @version         1.0
// @description     Community_Notification_System
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host      localhost:9080
// @BasePath  /api/v1
func main() {
	// 初始化 Gin
	router := gin.Default()

	//跨域的 Middleware
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.JWTAuthMiddleware())
	router.Use(middlewares.CookieMiddleware())
	//env 初始化
	configs.InitConfig()

	//DB 初始化
	database.InitDB()

	apiGroup := router.Group("/api")
	routers.RegisterRoutes(apiGroup)
	// 啟動伺服器

	// 將 Swagger UI 綁定在 /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	router.Run(port)
}
