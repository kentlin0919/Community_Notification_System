package message

import (
	"context"
	"fmt"
	"log"
	"net/http"

	message_model "Community_Notification_System/app/models/message"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/message"
	message_db "Community_Notification_System/database/Message_DB"
	"Community_Notification_System/pkg/firebase"

	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	utilsErr "Community_Notification_System/utils/errors"
)

const (
	messageTargetSelectedUsers = "selected_users"
	messageTargetCommunity     = "community"
	fcmStatusSent              = "sent"
	fcmStatusFailed            = "failed"
	fcmStatusSkipped           = "skipped"
)

// SendMessage 處理送入 message
// @Summary 傳送訊息
// @Description 管理者傳送社區訊息。系統會先為每位收件者建立 message_info 寄送紀錄，再嘗試透過 FCM 推播；Firebase 未初始化或部分推播失敗時仍會保留寄送紀錄並於回應中標示失敗數量。
// @Tags Message
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param message body message_model.SendMessageRequest true "訊息資料"
// @Success 200 {object} message_model.SendMessageResponse "訊息送出完成"
// @Failure 400 {object} model.Response400Error "輸入資料格式錯誤"
// @Failure 401 {object} model.Response401Error "未授權，缺少或無效的 JWT Token"
// @Failure 404 {object} model.Response404Error "找不到符合條件的收件者"
// @Failure 500 {object} model.Response500Error "伺服器內部錯誤"
// @Router /api/v1/messages/send [post]
// @Router /api/v1/sendmessage [post]
func (m *MessageController) SendMessage(ctx *gin.Context) {
	var req message_model.SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "無效的輸入資料")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	normalizeSendMessageRequest(&req)
	if err := validateSendMessageRequest(&req); err != nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, err.Error())
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	senderID, ok := getUserIDFromContext(ctx)
	if !ok {
		errorModel := model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "未授權")
		ctx.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	communityID, ok := getCommunityIDFromContext(ctx)
	if !ok {
		errorModel := model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "缺少社區上下文")
		ctx.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	recipientsResult := repository.FindMessageRecipientsRepository(communityID, &req)
	if recipientsResult.Statue.Error != nil {
		log.Printf("error finding message recipients: %v\n", recipientsResult.Statue.Error)
		errorModel := model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "查詢收件者失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	if len(recipientsResult.Result) == 0 {
		errorModel := model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrInvalidParams, "找不到符合條件的收件者")
		ctx.JSON(http.StatusNotFound, errorModel)
		return
	}

	batchID := uuid.NewString()
	recordsResult := repository.CreateMessageRecordsRepository(senderID, communityID, &req, recipientsResult.Result, batchID)
	if recordsResult.Statue.Error != nil {
		log.Printf("error creating message records: %v\n", recordsResult.Statue.Error)
		errorModel := model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "建立訊息紀錄失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	response := buildSendMessageResponse(batchID, recordsResult.Result.Records)
	for idx, record := range recordsResult.Result.Records {
		recipient := recipientsResult.Result[idx]
		result := sendMessageToRecipient(ctx, &req, &record, recipient.Fcmtoken)
		response.Data.FCMResults[idx] = result
		if result.Success {
			response.Data.SuccessCount++
		} else {
			response.Data.FailureCount++
		}
		repository.UpdateMessageFCMResultRepository(record.ID, result.FcmStatus, result.FcmMessageID, result.Error)
	}

	ctx.JSON(http.StatusOK, response)
}

func normalizeSendMessageRequest(req *message_model.SendMessageRequest) {
	if req.TargetType == "" {
		if len(req.RecipientUserIDs) > 0 || len(req.RecipientEmails) > 0 || len(req.Userselect) > 0 {
			req.TargetType = messageTargetSelectedUsers
		} else {
			req.TargetType = messageTargetCommunity
		}
	}
	if req.Subtitle == "" {
		req.Subtitle = req.Title
	}
}

func validateSendMessageRequest(req *message_model.SendMessageRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title 為必填")
	}
	if req.Body == "" {
		return fmt.Errorf("body 為必填")
	}
	if req.TargetType != messageTargetSelectedUsers && req.TargetType != messageTargetCommunity {
		return fmt.Errorf("target_type 僅支援 selected_users 或 community")
	}
	if req.TargetType == messageTargetSelectedUsers && len(req.RecipientUserIDs) == 0 && len(req.RecipientEmails) == 0 && len(req.Userselect) == 0 {
		return fmt.Errorf("selected_users 至少需要 recipient_user_ids 或 recipient_emails")
	}
	return nil
}

func buildSendMessageResponse(batchID string, records []message_db.MessageInfo) message_model.SendMessageResponse {
	data := message_model.SendMessageResponseData{
		MessageBatchID:    batchID,
		TargetCount:       len(records),
		CreatedMessageIDs: make([]string, len(records)),
		FCMResults:        make([]message_model.SendMessageResult, len(records)),
	}
	for idx, record := range records {
		data.CreatedMessageIDs[idx] = record.ID
		data.FCMResults[idx] = message_model.SendMessageResult{
			MessageID: record.ID,
			UserID:    record.UserID,
			Email:     record.Email,
		}
	}
	return message_model.SendMessageResponse{
		Message: "訊息發送完成",
		Data:    data,
	}
}

func sendMessageToRecipient(ctx *gin.Context, req *message_model.SendMessageRequest, record *message_db.MessageInfo, fcmToken string) message_model.SendMessageResult {
	result := message_model.SendMessageResult{
		MessageID: record.ID,
		UserID:    record.UserID,
		Email:     record.Email,
		FcmStatus: fcmStatusSkipped,
	}

	if firebase.FcmClient == nil {
		result.Error = "Firebase 推播服務尚未初始化"
		return result
	}
	if fcmToken == "" {
		result.Error = "收件者未註冊 FCM token"
		return result
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
		Data: map[string]string{
			"message_id":   record.ID,
			"batch_id":     record.BatchID,
			"category":     req.Metadata.Category,
			"entity_type":  req.Metadata.EntityType,
			"entity_id":    req.Metadata.EntityID,
			"click_action": req.Metadata.ClickAction,
			"community_id": fmt.Sprintf("%d", record.CommunityID),
			"notification": "message",
		},
		Token: fcmToken,
	}

	response, err := firebase.FcmClient.Send(context.Background(), message)
	if err != nil {
		log.Printf("error sending message to user %s: %v\n", record.UserID, err)
		result.FcmStatus = fcmStatusFailed
		result.Error = err.Error()
		return result
	}

	result.Success = true
	result.FcmStatus = fcmStatusSent
	result.FcmMessageID = response
	return result
}

func getCommunityIDFromContext(ctx *gin.Context) (uint64, bool) {
	value, exists := ctx.Get("community_id")
	if !exists {
		return 0, false
	}
	switch communityID := value.(type) {
	case uint64:
		return communityID, communityID > 0
	case uint:
		return uint64(communityID), communityID > 0
	case int:
		return uint64(communityID), communityID > 0
	case float64:
		return uint64(communityID), communityID > 0
	default:
		return 0, false
	}
}
