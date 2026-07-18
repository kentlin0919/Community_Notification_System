package v2

import (
	v1 "Community_Notification_System/app/controller/v1"
	"Community_Notification_System/middlewares"

	"github.com/gin-gonic/gin"
)

func V2PublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", v1.User().UserLogin)
	rg.POST("/register", v1.User().UserRegister)
}

func V2PrivateRoutes(rg *gin.RouterGroup) {
	// ── 一般登入使用者即可（無社區隔離限制） ──
	rg.POST("/deleteUser", v1.User().UserDelete)

	// ── 社區管理人員級（Staff+, Level 7，有社區資料隔離） ──
	staffScope := rg.Group("")
	staffScope.Use(middlewares.CommunityContextMiddleware(), middlewares.MinimalPermissionMiddleware(7))
	{
		staffScope.POST("/sendmessage", v1.Message().SendMessage)
		staffScope.POST("/messages/send", v1.Message().SendMessage)
	}
}
