package communityManager

import (
	accountModel "Community_Notification_System/app/models/account"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/community"
	userRepository "Community_Notification_System/app/repositories/user"
	communitydb "Community_Notification_System/database/Community_DB"
	userdb "Community_Notification_System/database/User_DB"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CommunityManager_Approve 超級管理員核可社區申請
// @Summary 核可社區申請
// @Description 僅允許 Super admin 操作，核可後建立正式社區與初始 admin 帳號。
// @Tags CommunityManager
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Success 200 {object} map[string]interface{} "核可成功，回傳社區與管理員資訊"
// @Failure 400 {object} model.Response400Error "請求錯誤或無此申請單"
// @Failure 403 {object} model.Response403Error "權限不足"
// @Failure 500 {object} model.Response500Error "核可流程失敗"
// @Security BearerAuth
// @Router /api/v1/community/register/{id}/approve [patch]
func (c *CommunityManagerController) CommunityManager_Approve(ctx *gin.Context) {
	idStr := ctx.Param("id")
	applicationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "無效的 Application ID"))
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

	// 構建將要建立的 Community Info
	communityInfo := &communitydb.CommunityInfo{
		PostalCode:     application.PostalCode,
		Municipality:   application.Municipality,
		District:       application.District,
		RoadName:       application.RoadName,
		LaneNumber:     application.LaneNumber,
		AlleyNumber:    application.AlleyNumber,
		Community_name: application.CommunityName,
		Address:        application.Address,
	}

	// 檢查是否已存在相同的社區
	checkCommunityStatue, err := checkCommunity(communityInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "檢查社區資料失敗"))
		return
	}
	if !checkCommunityStatue {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "此地址或社區參數已存在正式記錄"))
		return
	}

	// 檢查 admin_email 是否已註冊
	adminExist := userRepository.LoginRepository(&accountModel.User{Email: application.AdminEmail})
	if adminExist.Statue.Error == nil { // 表示找到資料了
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "申請單指定的管理員 Email 已經被註冊"))
		return
	}

	// 構建 admin UserInfo (Permission=2)
	adminUser := &userdb.UserInfo{
		ID:           uuid.New().String(),
		Email:        application.AdminEmail,
		Name:         application.AdminName,
		Password:     application.AdminPasswordHash, // 從 application 帶過來的 hash
		Registertime: time.Now(),
		PermissionId: 2, // Admin 預設權限
	}

	// 使用 Transaction 進行核可
	err = repository.ApproveApplicationTransactionRepository(&application, communityInfo, adminUser, currentUserInfo.Result.ID)
	if err != nil {
		fmt.Printf("核可失敗: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "核可流程失敗: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "社區申請核可成功",
		"data": gin.H{
			"community_id":   communityInfo.Community_id,
			"community_name": communityInfo.Community_name,
			"admin_email":    adminUser.Email,
			"status":         "approved",
		},
	})
}
