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

// @title           My Swagger Demo
// @version         1.0
// @description     這是一個範例服務，示範如何建立 Swagger 文件。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9080
// @BasePath  /api/v1
func main() {
	// 初始化 Gin
	router := gin.Default()
	//跨域的 Middleware
	router.Use(middlewares.CORSMiddleware())

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
