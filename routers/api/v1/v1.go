package v1

import (
	v1 "Community_Notification_System/app/controller/v1"

	"github.com/gin-gonic/gin"
)

func V1Routes(rg *gin.RouterGroup) {

	rg.POST("/login", v1.User().UserLogin)
	/// 處理註冊請求
	rg.POST("/register", v1.User().UserRegister)

}
