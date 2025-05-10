package user

import (
	accountModel "Community_Notification_System/app/models/account"
	models "Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserLogin 處理刪除使用者
// @Summary 使用者刪除
// @Description 使用者提供帳號與密碼後刪除使用者
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param login body accountModel.User true "使用者資料（Email & Password ＆ Platform）"
// @Success 200 {object} models.RequestMessage "登入成功，返回 JWT Token 和成功訊息"
// @Failure 400 {object} models.ErrorRequest "無效的輸入資料"
// @Failure 401 {object} models.ErrorRequest "密碼錯誤"
// @Failure 404 {object} model.ErrorRequest "使用者不存在"
// @Failure 500 {object} models.ErrorRequest "系統錯誤或 JWT 簽發失敗"
// @Router /api/v1/deleteUser [post]
func (u *UserController) UserDelete(ctx *gin.Context) {
	var UserDeleteModel accountModel.User

	// 綁定 JSON 資料
	if err := ctx.ShouldBindJSON(&UserDeleteModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	//搜尋資料庫是否有此用戶
	result := repository.LoginRepository(&UserDeleteModel)

	if result.Statue.Error != nil {
		if result.Statue.Error == gorm.ErrRecordNotFound {
			var errorModel models.ErrorRequest
			errorModel.Error = "User Not Found"
			ctx.JSON(http.StatusNotFound, errorModel)
			return
		}

		var errorModel models.ErrorRequest
		errorModel.Error = "System Error"
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	deleteResult := repository.UserDeleteRepository(&result.Result)

	if deleteResult.Statue.Error != nil {
		var errorModel models.ErrorRequest
		errorModel.Error = "Delete Error"
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	var RequestMessage models.RequestMessage
	RequestMessage.Message = "Delete Sucessful"
	ctx.JSON(http.StatusAccepted, RequestMessage)
}
