package v2

import (
	v1 "Community_Notification_System/app/controller/v1"

	"github.com/gin-gonic/gin"
)

func V2PublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", v1.User().UserLogin)
	rg.POST("/register", v1.User().UserRegister)
}

func V2PrivateRoutes(rg *gin.RouterGroup) {
	rg.POST("/deleteUser", v1.User().UserDelete)
	rg.POST("/sendmessage", v1.Message().SendMessage)
}
