package message

import (
	message_model "Community_Notification_System/app/models/message"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/message"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	utilsErr "Community_Notification_System/utils/errors"
)

// GetMessageList 取得當前使用者通知列表
// @Summary 取得通知列表
// @Description 依登入使用者取得通知列表，支援分頁與已讀篩選。
// @Tags Message
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "頁碼（預設 1）"
// @Param page_size query int false "每頁筆數（預設 20，最大 100）"
// @Param is_read query bool false "已讀篩選"
// @Success 200 {object} message_model.MessageListResponse "成功取得通知列表"
// @Failure 400 {object} model.Response400Error "請求參數錯誤"
// @Failure 401 {object} model.Response401Error "未授權"
// @Failure 500 {object} model.Response500Error "系統錯誤"
// @Router /api/v1/messages [get]
func (m *MessageController) GetMessageList(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		errorModel := model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "未授權")
		ctx.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	var query message_model.MessageListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "請求參數錯誤")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	repoResult := repository.GetUserMessagesRepository(userID, &query)
	if repoResult.Statue.Error != nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "取得通知失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	response := message_model.MessageListResponse{
		Total:    repoResult.Result.Total,
		Messages: make([]message_model.MessageSummary, len(repoResult.Result.Items)),
	}

	for idx, item := range repoResult.Result.Items {
		response.Messages[idx] = message_model.MessageSummary{
			ID:         item.ID,
			Title:      item.Title,
			Subtile:    item.Subtile,
			Detail:     item.Detail,
			IsRead:     item.IsRead,
			CreateTime: item.CreateTime.Format(time.RFC3339),
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// MarkMessageRead 標記單筆通知為已讀
// @Summary 標記通知為已讀
// @Description 將登入使用者的指定通知標記為已讀。
// @Tags Message
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "通知 ID"
// @Success 200 {object} model.RequestMessage "成功標記已讀"
// @Failure 401 {object} model.Response401Error "未授權"
// @Failure 404 {object} model.Response404Error "通知不存在"
// @Failure 500 {object} model.Response500Error "系統錯誤"
// @Router /api/v1/messages/{id}/read [patch]
func (m *MessageController) MarkMessageRead(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		errorModel := model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "未授權")
		ctx.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	messageID := ctx.Param("id")
	repoResult := repository.MarkMessageReadRepository(userID, messageID)
	if repoResult.Statue.Error != nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "標記已讀失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	if !repoResult.Result {
		errorModel := model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "通知不存在")
		ctx.JSON(http.StatusNotFound, errorModel)
		return
	}

	ctx.JSON(http.StatusOK, model.RequestMessage{Message: "標記已讀成功"})
}

// MarkAllMessagesRead 標記使用者全部通知為已讀
// @Summary 全部標記已讀
// @Description 將登入使用者所有未讀通知標記為已讀。
// @Tags Message
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.RequestMessage "成功全部標記已讀"
// @Failure 401 {object} model.Response401Error "未授權"
// @Failure 500 {object} model.Response500Error "系統錯誤"
// @Router /api/v1/messages/read-all [patch]
func (m *MessageController) MarkAllMessagesRead(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		errorModel := model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "未授權")
		ctx.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	repoResult := repository.MarkAllMessagesReadRepository(userID)
	if repoResult.Statue.Error != nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "全部標記已讀失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	ctx.JSON(http.StatusOK, model.RequestMessage{Message: "全部標記已讀成功"})
}

func getUserIDFromContext(ctx *gin.Context) (string, bool) {
	value, exists := ctx.Get("user_id")
	if !exists {
		return "", false
	}
	userID, ok := value.(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}
