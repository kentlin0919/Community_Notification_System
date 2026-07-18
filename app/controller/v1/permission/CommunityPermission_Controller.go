package permission

import (
	"Community_Notification_System/app/models/model"
	permissionModel "Community_Notification_System/app/models/permission"
	repository "Community_Notification_System/app/repositories/permission"
	permissiondb "Community_Notification_System/database/Permission_DB"
	"net/http"

	"github.com/gin-gonic/gin"

	utilsErr "Community_Notification_System/utils/errors"
)

type CommunityPermissionController struct{}

func NewCommunityPermissionController() *CommunityPermissionController {
	return &CommunityPermissionController{}
}

// UpdateCommunityPermissionProfile 更新或新增社區角色的顯示名稱
// @Summary 管理社區專屬角色名稱
// @Description 社區 Admin (或 Super admin) 可修改 PermissionID = 3~7 的顯示名稱。用於客製化不同社區對同一權限等級的稱呼。
// @Tags Permission Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param community_id header string true "社區 ID"
// @Param body body permissionModel.UpdatePermissionProfileRequest true "更新請求內容"
// @Success 200 {object} model.RequestMessage "權限名稱更新成功"
// @Failure 400 {object} model.Response400Error "請求參數錯誤或 ID 範圍不符"
// @Failure 401 {object} model.Response401Error "認證失敗或無法取得社區資訊"
// @Failure 403 {object} model.Response403Error "權限不足"
// @Failure 500 {object} model.Response500Error "系統更新失敗"
// @Router /api/v1/admin/permissions/profile [put]
func (c *CommunityPermissionController) UpdateCommunityPermissionProfile(ctx *gin.Context) {
	// 從 CommunityContextMiddleware 及 Jwt 中拿到隔離資訊
	communityID, exists := ctx.Get("community_id")
	userID, userExists := ctx.Get("user_id")

	if !exists || !userExists {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}

	var req permissionModel.UpdatePermissionProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid input: "+err.Error()))
		return
	}

	// Permission 範圍限制已在 Model Bound 或是這邊再做一次確保
	if req.PermissionID < 3 || req.PermissionID > 7 {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "僅允許修改 Permission ID 介於 3 至 7 的自訂角色"))
		return
	}

	cID := communityID.(uint64)
	uID := userID.(string)

	profile := &permissiondb.CommunityPermissionProfile{
		CommunityID:  cID,
		PermissionID: req.PermissionID,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		UpdatedBy:    uID,
	}

	res := repository.UpsertPermissionProfileRepository(profile)
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "更新權限設定失敗"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "權限名稱更新成功",
		"data":    res.Result,
	})
}

// GetCommunityPermissionProfiles 查詢特定社區自訂角色與預設角色的對照
// @Summary 查詢社區角色列表
// @Description 取得社區客製化的 3~7 等級角色名稱與描述。
// @Tags Permission Management
// @Produce json
// @Security BearerAuth
// @Param community_id header string true "社區 ID"
// @Success 200 {object} model.RequestMessage "成功回傳角色對照表"
// @Failure 401 {object} model.Response401Error "認證失敗"
// @Failure 500 {object} model.Response500Error "查詢失敗"
// @Router /api/v1/permissions/profile [get]
func (c *CommunityPermissionController) GetCommunityPermissionProfiles(ctx *gin.Context) {
	communityID, exists := ctx.Get("community_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}

	cID := communityID.(uint64)
	res := repository.GetCommunityPermissionProfilesRepository(cID)
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "查詢失敗"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": res.Result,
	})
}
