package routers

import (
	"Community_Notification_System/middlewares"
	"github.com/gin-gonic/gin"

	//引入v1 的router
	v1 "Community_Notification_System/routers/api/v1"
	v2 "Community_Notification_System/routers/api/v2"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	// v1
	v1Group := rg.Group("/v1")
	v1.V1PublicRoutes(v1Group)
	v1PrivateGroup := rg.Group("/v1")
	v1PrivateGroup.Use(middlewares.JWTAuthMiddleware())
	v1.V1PrivateRoutes(v1PrivateGroup)

	// v2
	v2Group := rg.Group("/v2")
	v2.V2PublicRoutes(v2Group)
	v2PrivateGroup := rg.Group("/v2")
	v2PrivateGroup.Use(middlewares.JWTAuthMiddleware())
	v2.V2PrivateRoutes(v2PrivateGroup)
}
