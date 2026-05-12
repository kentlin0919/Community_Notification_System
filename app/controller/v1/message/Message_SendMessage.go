package message

import (
	message_model "Community_Notification_System/app/models/message"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/user"
	"Community_Notification_System/pkg/firebase"
	"context"
	"log"
	"net/http"
	"Community_Notification_System/utils/errors"

	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
)

// SendMessage 處理送入 message
// @Summary 傳送訊息
// @Description 使用者傳送訊息，需要 JWT Token 驗證
// @Tags Message
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param message body message_model.FCMNotificationRequest true "訊息資料"
// @Success 200 {object} message_model.MessageRequest "訊息送出成功"
// @Failure 400 {object} model.Response400Error "輸入資料格式錯誤"
// @Failure 401 {object} model.Response401Error "未授權，缺少或無效的 JWT Token"
// @Failure 500 {object} model.Response500Error "伺服器內部錯誤"
// @Router /api/v1/sendmessage [post]
func (m *MessageController) SendMessage(ctx *gin.Context) {

	var req message_model.FCMNotificationRequest

	// 綁定 JSON 資料並驗證輸入格式
	// 使用 ShouldBindJSON 可以自動驗證 JSON 格式是否符合結構體定義
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusBadRequest, errors.ErrInvalidParams, "無效的輸入資料")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	_ = repository.UserInfoListRepository(req.Userselect) // TODO: refactor to use message domain interface

	if firebase.FcmClient == nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusServiceUnavailable, errors.ErrInternal, "Firebase 推播服務尚未初始化")
		ctx.JSON(http.StatusServiceUnavailable, errorModel)
		return
	}

	// The message to send.
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
		Token: req.DeviceToken, // This is the device token obtained from the client app.
	}

	// Send the message.
	response, err := firebase.FcmClient.Send(context.Background(), message)
	if err != nil {
		log.Printf("error sending message: %v\n", err)
		errorModel := model.NewErrorResponse(ctx, http.StatusInternalServerError, errors.ErrInternal, "Failed to send message")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	log.Printf("Successfully sent message: %s\n", response)
	ctx.JSON(200, gin.H{"message": "Successfully sent message", "response": response})

}
