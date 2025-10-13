package v1

import (
	v1 "Community_Notification_System/app/controller/v1"

	"github.com/gin-gonic/gin"
)

func V1Routes(rg *gin.RouterGroup) {

	// 處理登入請求
	rg.POST("/login", v1.User().UserLogin)
	/// 處理註冊請求
	rg.POST("/register", v1.User().UserRegister)

	/// 刪除使用者
	rg.POST("/deleteUser", v1.User().UserDelete)

	//處理送通知
	rg.POST("/sendmessage", v1.Message().SendMessage)

	// 取得社區列表
	rg.GET("/community/getlist", v1.CommunityManager().CommunityManager_GetList)

	// 新增社區
	rg.POST("/community/register", v1.CommunityManager().CommunityManager_Register)

	// 取得平台列表
	rg.GET("/platform/getlist", v1.Platform().Platform_GetList)

}
