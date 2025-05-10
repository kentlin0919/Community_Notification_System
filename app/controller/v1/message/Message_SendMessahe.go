package message

import (
	message_model "Community_Notification_System/app/models/message"
	"Community_Notification_System/app/models/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendMessage 處理送入 message
// @Summary 傳送訊息
// @Description 使用者傳送訊息，需要 JWT Token 驗證
// @Tags Message
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param message body message_model.MessageData true "訊息資料"
// @Success 200 {object} message_model.MessageRequest "訊息送出成功"
// @Failure 400 {object} model.ErrorRequest "輸入資料格式錯誤"
// @Failure 401 {object} model.ErrorRequest "未授權，缺少或無效的 JWT Token"
// @Failure 500 {object} model.ErrorRequest "伺服器內部錯誤"
// @Router /api/v1/sendmessage [post]
func (m *MessageController) SendMessage(ctx *gin.Context) {

	var message_model message_model.MessageData

	// 綁定 JSON 資料並驗證輸入格式
	// 使用 ShouldBindJSON 可以自動驗證 JSON 格式是否符合結構體定義
	if err := ctx.ShouldBindJSON(&message_model); err != nil {
		var errorModel model.ErrorRequest
		errorModel.Error = "無效的輸入資料"
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorModel})
		return
	}

}
