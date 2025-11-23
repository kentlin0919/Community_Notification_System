package message

import (
	message_model "Community_Notification_System/app/models/message"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/user"
	"Community_Notification_System/pkg/firebase"
	"context"
	"fmt"
	"log"
	"net/http"

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
// @Failure 400 {object} model.ErrorRequest "輸入資料格式錯誤"
// @Failure 401 {object} model.ErrorRequest "未授權，缺少或無效的 JWT Token"
// @Failure 500 {object} model.ErrorRequest "伺服器內部錯誤"
// @Router /api/v1/sendmessage [post]
func (m *MessageController) SendMessage(ctx *gin.Context) {

	var req message_model.FCMNotificationRequest

	// 綁定 JSON 資料並驗證輸入格式
	// 使用 ShouldBindJSON 可以自動驗證 JSON 格式是否符合結構體定義
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "無效的輸入資料")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	UserInfoList := repository.UserInfoListRepository(req.Userselect)

	fmt.Print(UserInfoList)

	// The message to send.
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
		Token: req.DeviceToken, // This is the device token obtained from the client app.
	}

	// Check if Firebase client is initialized
	if firebase.FcmClient == nil {
		log.Println("Firebase client is not initialized")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Firebase service not available"})
		return
	}

	// Send the message.
	response, err := firebase.FcmClient.Send(context.Background(), message)
	if err != nil {
		log.Printf("error sending message: %v\n", err)
		ctx.JSON(500, gin.H{"error": "Failed to send message"})
		return
	}

	log.Printf("Successfully sent message: %s\n", response)
	ctx.JSON(200, gin.H{"message": "Successfully sent message", "response": response})

}
