package facility

import (
	facilityModel "Community_Notification_System/app/models/facility"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/facility"
	"Community_Notification_System/database"
	facilitydb "Community_Notification_System/database/Facility_DB"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	utilsErr "Community_Notification_System/utils/errors"
)

type FacilityController struct{}

func NewFacilityController() *FacilityController {
	return &FacilityController{}
}

// CreateFacility 新增預約設施
// @Summary 新增預約設施
// @Description 社區管理員新增公用設施與預設預約條件。系統會同時建立一筆預設的設施規則。
// @Tags Facility Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body facilityModel.CreateFacilityRequest true "設施基本資料"
// @Success 201 {object} model.RequestMessage "設施新增成功"
// @Failure 400 {object} model.Response400Error "請求參數錯誤"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊或認證失敗"
// @Failure 403 {object} model.Response403Error "權限不足，僅限管理員"
// @Failure 409 {object} model.Response409Error "該社區已存在相同名稱的設施"
// @Failure 500 {object} model.Response500Error "系統建立設施失敗"
// @Router /api/v1/admin/facilities [post]
func (c *FacilityController) CreateFacility(ctx *gin.Context) {
	communityID, exists := ctx.Get("community_id")
	userID, userExists := ctx.Get("user_id")

	if !exists || !userExists {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}

	var req facilityModel.CreateFacilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid input: "+err.Error()))
		return
	}

	cID := communityID.(uint64)

	// API 防呆：只能替自己的社區建設施 (這個防線在 CommunityContextMiddleware 也有一道)
	if req.CommunityID != cID {
		ctx.JSON(http.StatusForbidden, model.NewErrorResponse(ctx, http.StatusForbidden, utilsErr.ErrForbidden, "你只能管理所屬社區之設施"))
		return
	}

	status := "draft"
	if req.Status != "" {
		status = req.Status
	}

	newFacility := &facilitydb.FacilityInfo{
		CommunityID:  cID,
		Name:         req.Name,
		FacilityType: req.FacilityType,
		Location:     req.Location,
		Description:  req.Description,
		Status:       status,
		CoverImage:   req.CoverImage,
		CreatedBy:    userID.(string),
	}

	// 呼叫 Repository 新增
	res := repository.CreateFacilityWithRuleTransactionRepository(newFacility)
	if res.Statue.Error != nil {
		// 判斷是否為重複命名
		if res.Statue.Error.Error() == "該社區已存在相同名稱的設施" {
			ctx.JSON(http.StatusConflict, model.NewErrorResponse(ctx, http.StatusConflict, utilsErr.ErrConflict, res.Statue.Error.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "建立設施失敗"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "設施新增成功",
		"data": gin.H{
			"id":            res.Result.ID,
			"community_id":  res.Result.CommunityID,
			"name":          res.Result.Name,
			"facility_type": res.Result.FacilityType,
			"status":        res.Result.Status,
		},
	})
}

func permissionIDFromContext(ctx *gin.Context) int {
	permissionID, exists := ctx.Get("permission_id")
	if !exists {
		return 3
	}

	switch value := permissionID.(type) {
	case int:
		return value
	case float64:
		return int(value)
	case string:
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return 3
}

func communityIDFromContext(ctx *gin.Context) (uint64, bool) {
	communityID, exists := ctx.Get("community_id")
	if !exists {
		return 0, false
	}

	switch value := communityID.(type) {
	case uint64:
		return value, true
	case uint:
		return uint64(value), true
	case int:
		return uint64(value), true
	case float64:
		return uint64(value), true
	case string:
		parsed, err := strconv.ParseUint(value, 10, 64)
		return parsed, err == nil
	}
	return 0, false
}

// GetFacilityList 取得設施列表
// @Summary 取得設施列表
// @Description 取得設施列表。可經由 Query 參數指定社區 ID；若未指定，則從 JWT Token 中自動提取。住戶只會看到 active 設施；管理員可看到同社區全部未刪除設施。
// @Tags Facility Management
// @Produce json
// @Security BearerAuth
// @Param community_id query uint64 false "社區 ID"
// @Success 200 {object} facilityModel.FacilityListResponse "設施列表"
// @Failure 400 {object} model.Response400Error "請求參數錯誤"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 403 {object} model.Response403Error "權限不足"
// @Failure 500 {object} model.Response500Error "查詢設施列表失敗"
// @Router /api/v1/facilities [get]
func (c *FacilityController) GetFacilityList(ctx *gin.Context) {
	var communityID uint64

	// 1. 優先從 Query 參數中取得 community_id
	if queryCommunityIDStr := ctx.Query("community_id"); queryCommunityIDStr != "" {
		parsed, err := strconv.ParseUint(queryCommunityIDStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "無效的社區 ID 格式"))
			return
		}
		communityID = parsed
	} else {
		// 2. Fallback 到 JWT Context 中的 community_id
		jwtCommunityID, ok := communityIDFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "請提供社區 ID 或重新登入以獲取所屬社區資訊"))
			return
		}
		communityID = jwtCommunityID
	}

	// 3. 進行安全性與權限校驗
	permissionID := permissionIDFromContext(ctx)
	if permissionID > 1 {
		// 非超級管理員：校驗其所查詢之 community_id 是否與 JWT 所綁定之一致
		jwtCommunityID, ok := communityIDFromContext(ctx)
		if ok && communityID != jwtCommunityID {
			ctx.JSON(http.StatusForbidden, model.NewErrorResponse(ctx, http.StatusForbidden, utilsErr.ErrForbidden, "您無權操作或查詢其他社區的資料"))
			return
		}
	}

	onlyActive := permissionID > 2
	res := repository.GetFacilityListRepository(communityID, onlyActive)
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "取得設施列表失敗"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": res.Result})
}

// GetFacilityDetail 取得單一設施詳情
// @Summary 取得單一設施詳情
// @Description 依 JWT 社區取得單一設施詳情。
// @Tags Facility Management
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "設施 ID"
// @Success 200 {object} facilityModel.FacilityDetailResponse "設施詳情"
// @Failure 400 {object} model.Response400Error "無效的設施 ID"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 404 {object} model.Response404Error "找不到設施"
// @Router /api/v1/facilities/{id} [get]
func (c *FacilityController) GetFacilityDetail(ctx *gin.Context) {
	communityID, ok := communityIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid facility ID"))
		return
	}

	res := repository.GetFacilityDetailRepository(id, communityID)
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusNotFound, model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "找不到該設施"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": res.Result})
}

// UpdateFacility 更新設施
// @Summary 更新設施
// @Description 社區管理員更新設施基本資訊與預約規則。
// @Tags Facility Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "設施 ID"
// @Param body body facilityModel.UpdateFacilityRequest true "設施更新資料"
// @Success 200 {object} model.RequestMessage "設施更新成功"
// @Failure 400 {object} model.Response400Error "無效輸入"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 403 {object} model.Response403Error "權限不足"
// @Failure 404 {object} model.Response404Error "找不到設施"
// @Failure 409 {object} model.Response409Error "設施名稱衝突"
// @Router /api/v1/admin/facilities/{id} [put]
func (c *FacilityController) UpdateFacility(ctx *gin.Context) {
	if permissionIDFromContext(ctx) > 2 {
		ctx.JSON(http.StatusForbidden, model.NewErrorResponse(ctx, http.StatusForbidden, utilsErr.ErrForbidden, "權限不足，僅限管理員"))
		return
	}
	communityID, ok := communityIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid facility ID"))
		return
	}

	var req facilityModel.UpdateFacilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid input: "+err.Error()))
		return
	}

	detail := repository.GetFacilityDetailRepository(id, communityID)
	if detail.Statue.Error != nil {
		ctx.JSON(http.StatusNotFound, model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "找不到該設施"))
		return
	}

	facilityInfo := detail.Result
	facilityInfo.Name = req.Name
	facilityInfo.FacilityType = req.FacilityType
	facilityInfo.Location = req.Location
	facilityInfo.Description = req.Description
	facilityInfo.CoverImage = req.CoverImage
	if req.Status != "" {
		facilityInfo.Status = req.Status
	}

	var rule facilitydb.FacilityRule
	rulePointer := (*facilitydb.FacilityRule)(nil)
	if err := database.DB.Where("facility_id = ?", id).First(&rule).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "查詢設施規則失敗"))
		return
	}
	rulePointer = &rule
	if req.SlotMinutes != nil {
		rule.SlotMinutes = *req.SlotMinutes
	}
	if req.MaxAdvanceDays != nil {
		rule.MaxAdvanceDays = *req.MaxAdvanceDays
	}
	if req.CancelBeforeHours != nil {
		rule.CancelBeforeHours = *req.CancelBeforeHours
	}
	if req.AutoApprove != nil {
		rule.AutoApprove = *req.AutoApprove
	}
	if req.ReservationRuleOpen != nil {
		rule.IsActive = *req.ReservationRuleOpen
	}

	updateRes := repository.UpdateFacilityWithRuleRepository(&facilityInfo, rulePointer)
	if updateRes.Statue.Error != nil {
		if errors.Is(updateRes.Statue.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "找不到該設施"))
			return
		}
		ctx.JSON(http.StatusConflict, model.NewErrorResponse(ctx, http.StatusConflict, utilsErr.ErrConflict, updateRes.Statue.Error.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "設施更新成功", "data": updateRes.Result})
}

// DeleteFacility 刪除設施
// @Summary 軟刪除設施
// @Description 社區管理員軟刪除設施，既有預約歷史保留。
// @Tags Facility Management
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "設施 ID"
// @Success 200 {object} model.RequestMessage "設施刪除成功"
// @Failure 400 {object} model.Response400Error "無效的設施 ID"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 403 {object} model.Response403Error "權限不足"
// @Failure 404 {object} model.Response404Error "找不到設施"
// @Router /api/v1/admin/facilities/{id} [delete]
func (c *FacilityController) DeleteFacility(ctx *gin.Context) {
	if permissionIDFromContext(ctx) > 2 {
		ctx.JSON(http.StatusForbidden, model.NewErrorResponse(ctx, http.StatusForbidden, utilsErr.ErrForbidden, "權限不足，僅限管理員"))
		return
	}
	communityID, ok := communityIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid facility ID"))
		return
	}

	if err := repository.DeleteFacilityRepository(id, communityID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "找不到該設施"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "刪除設施失敗"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "設施刪除成功"})
}
