package user

import (
	accountModel "Community_Notification_System/app/models/account"
	models "Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	utilsErr "Community_Notification_System/utils/errors"
)

// UserDelete 處理刪除使用者
// @Summary 刪除使用者帳號
// @Description 使用者提供帳號與密碼後進行帳號刪除。此操作為實體刪除且不可逆，執行前請務必確認。
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param login body accountModel.User true "刪除驗證資料（Email, Password, Platform）"
// @Success 202 {object} models.RequestMessage "帳號刪除成功"
// @Failure 400 {object} models.Response400Error "請求格式錯誤或無效輸入"
// @Failure 401 {object} models.Response401Error "認證失敗或密碼錯誤"
// @Failure 404 {object} models.Response404Error "找不到該使用者"
// @Failure 500 {object} models.Response500Error "系統執行刪除時發生錯誤"
// @Router /api/v1/deleteUser [post]
func (u *UserController) UserDelete(ctx *gin.Context) {
	var UserDeleteModel accountModel.User

	// 綁定 JSON 資料
	if err := ctx.ShouldBindJSON(&UserDeleteModel); err != nil {
		errorModel := models.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid input")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	//搜尋資料庫是否有此用戶
	result := repository.LoginRepository(&UserDeleteModel)

	if result.Statue.Error != nil {
		if result.Statue.Error == gorm.ErrRecordNotFound {
			errorModel := models.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "User Not Found")
			ctx.JSON(http.StatusNotFound, errorModel)
			return
		}

		errorModel := models.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "System Error")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	deleteResult := repository.UserDeleteRepository(&result.Result)

	if deleteResult.Statue.Error != nil {
		errorModel := models.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "Delete Error")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	var RequestMessage models.RequestMessage
	RequestMessage.Message = "Delete Sucessful"
	ctx.JSON(http.StatusAccepted, RequestMessage)
}
