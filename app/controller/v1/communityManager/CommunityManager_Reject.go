package communityManager

import (
	accountModel "Community_Notification_System/app/models/account"
	communityModel "Community_Notification_System/app/models/community"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/community"
	userRepository "Community_Notification_System/app/repositories/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Failure 400 {object} model.ErrorRequest "請求錯誤或無此申請單"
// @Failure 403 {object} model.ErrorRequest "權限不足"
// @Failure 500 {object} model.ErrorRequest "駁回流程失敗"
// @Security BearerAuth
// @Router /api/v1/community/register/{id}/reject [patch]
func (c *CommunityManagerController) CommunityManager_Reject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	applicationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "無效的 Application ID"))
		return
	}

	var req communityModel.RejectApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid input"))
		return
	}

	// 取得 jwt 解析後的 email
	currentUserEmail := ctx.GetString("username")
	if currentUserEmail == "" {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorRequest(http.StatusUnauthorized, "無法取得使用者登入資訊"))
		return
	}

	// 取得登入者的資訊
	currentUserInfo := userRepository.LoginRepository(&accountModel.User{Email: currentUserEmail})
	if currentUserInfo.Statue.Error != nil {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorRequest(http.StatusUnauthorized, "無法驗證登入者身分"))
		return
	}

	// 檢查是否為 Super admin (PermissionID == 1)
	if currentUserInfo.Result.PermissionId != 1 {
		ctx.JSON(http.StatusForbidden, model.NewErrorRequest(http.StatusForbidden, "僅 Super admin 可審核社區申請"))
		return
	}

	// 取得申請單
	appRes := repository.GetApplicationByIDRepository(applicationID)
	if appRes.Statue.Error != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "找不到該筆申請"))
		return
	}

	application := appRes.Result
	if application.Status != "pending" {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "此申請單已不在待審核狀態"))
		return
	}

	updateRes := repository.UpdateApplicationStatusRepository(applicationID, "rejected", currentUserInfo.Result.ID, req.RejectReason)
	if updateRes.Statue.Error != nil {
		fmt.Printf("駁回失敗: %v\n", updateRes.Statue.Error)
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "駁回狀態更新失敗"))
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
