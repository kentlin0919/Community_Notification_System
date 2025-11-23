package apisync

import (
	"log"

	"Community_Notification_System/database"
	"Community_Notification_System/database/ApiRoute_DB"

	"github.com/gin-gonic/gin"
)

// SyncApiRoutes 將 Gin 路由同步到資料庫，供權限控管使用。
func SyncApiRoutes(router *gin.Engine) {
	routes := router.Routes()
	log.Printf("發現 %d 個路由，開始同步至資料庫...", len(routes))

	for _, route := range routes {
		// 僅同步 /api/ 開頭的路由
		if len(route.Path) < 5 || route.Path[:5] != "/api/" {
			continue
		}

		var existingRoute ApiRoute_DB.ApiRoute
		if err := database.DB.Where("path = ? AND method = ?", route.Path, route.Method).First(&existingRoute).Error; err != nil {
			newRoute := ApiRoute_DB.ApiRoute{
				Path:        route.Path,
				Method:      route.Method,
				Description: route.Handler,
			}
			if createErr := database.DB.Create(&newRoute).Error; createErr != nil {
				log.Printf("無法新增路由 %s %s: %v", route.Method, route.Path, createErr)
			} else {
				log.Printf("新增路由: %s %s", route.Method, route.Path)
			}
		}
	}
	log.Println("API 路由同步完成。")
}
