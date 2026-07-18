package communityManager

import (
	communityModel "Community_Notification_System/app/models/community"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/community"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	utilsErr "Community_Notification_System/utils/errors"
)

// CommunityManager_Reject 超級管理員駁回社區申請
// @Summary 駁回社區申請
// @Description 僅允許 Super admin 操作，駁回申請並提供原因。
// @Tags CommunityManager
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Param body body communityModel.RejectApplicationRequest true "駁回原因"
// @Success 200 {object} map[string]interface{} "已駁回申請"
// @Failure 400 {object} model.Response400Error "請求錯誤或無此申請單"
// @Failure 403 {object} model.Response403Error "權限不足"
// @Failure 500 {object} model.Response500Error "駁回流程失敗"
// @Security BearerAuth
// @Router /api/v1/community/register/{id}/reject [patch]
func (c *CommunityManagerController) CommunityManager_Reject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	applicationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "無效的 Application ID"))
		return
	}

	var req communityModel.RejectApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid input"))
		return
	}

	// 從 JWT Middleware 注入的 Context 取得審核者 user_id
	rejectorUserID := ctx.GetString("user_id")
	if rejectorUserID == "" {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法取得使用者登入資訊"))
		return
	}

	// 取得申請單
	appRes := repository.GetApplicationByIDRepository(applicationID)
	if appRes.Statue.Error != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "找不到該筆申請"))
		return
	}

	application := appRes.Result
	if application.Status != "pending" {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "此申請單已不在待審核狀態"))
		return
	}

	updateRes := repository.UpdateApplicationStatusRepository(applicationID, "rejected", rejectorUserID, req.RejectReason)
	if updateRes.Statue.Error != nil {
		fmt.Printf("駁回失敗: %v\n", updateRes.Statue.Error)
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "駁回狀態更新失敗"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "社區申請已駁回",
		"data": gin.H{
			"application_id": updateRes.Result.ID,
			"status":         updateRes.Result.Status,
		},
	})
}
