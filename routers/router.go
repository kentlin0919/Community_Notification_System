package routers

import (
	"github.com/gin-gonic/gin"

	//引入v1 的router
	v1 "Community_Notification_System/routers/api/v1"
	v2 "Community_Notification_System/routers/api/v2"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	// 創建 v1 路由組
	v1Group := rg.Group("/v1")
	v1.V1Routes(v1Group) // 引入 v1.go 的路由

	v2Group := rg.Group("/v2")

	v2.V2Routes(v2Group)
}
