package user

import (
	accountModel "Community_Notification_System/app/models/account"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *UserController) UserDelete(ctx *gin.Context) {
	var registerModel accountModel.Register

	// 綁定 JSON 資料
	if err := ctx.ShouldBindJSON(&registerModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

}
