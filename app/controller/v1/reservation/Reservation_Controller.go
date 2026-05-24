package reservation

import (
	"Community_Notification_System/app/models/model"
	reservationModel "Community_Notification_System/app/models/reservation"
	repository "Community_Notification_System/app/repositories/reservation"
	"Community_Notification_System/database"
	facilitydb "Community_Notification_System/database/Facility_DB"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"


	utilsErr "Community_Notification_System/utils/errors"
)

type ReservationController struct{}

func NewReservationController() *ReservationController {
	return &ReservationController{}
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

func userIDFromContext(ctx *gin.Context) (string, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return "", false
	}
	userIDString, ok := userID.(string)
	return userIDString, ok && userIDString != ""
}

// GetReservationList 取得預約列表
// @Summary 取得預約列表
// @Description 依 JWT 社區取得預約列表。住戶只會看到自己的預約；管理員可看到同社區全部未刪除預約。
// @Tags Reservation
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "預約列表"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 500 {object} model.Response500Error "查詢預約列表失敗"
// @Router /api/v1/reservations [get]
func (c *ReservationController) GetReservationList(ctx *gin.Context) {
	communityID, ok := communityIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}
	userID, ok := userIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法取得使用者資訊，請重新登入"))
		return
	}

	filterByUserID := permissionIDFromContext(ctx) > 2
	result := repository.GetReservationListRepository(communityID, userID, filterByUserID)
	if result.Statue.Error != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "取得預約列表失敗"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result.Result})
}

// GetReservationDetail 取得單一預約詳情
// @Summary 取得預約詳情
// @Description 依 JWT 社區與角色範圍取得單一預約詳情。
// @Tags Reservation
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "預約 ID"
// @Success 200 {object} map[string]interface{} "預約詳情"
// @Failure 400 {object} model.Response400Error "無效的預約 ID"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 404 {object} model.Response404Error "找不到預約"
// @Router /api/v1/reservations/{id} [get]
func (c *ReservationController) GetReservationDetail(ctx *gin.Context) {
	communityID, ok := communityIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}
	userID, ok := userIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, "無法取得使用者資訊，請重新登入"))
		return
	}
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid reservation ID"))
		return
	}

	filterByUserID := permissionIDFromContext(ctx) > 2
	result := repository.GetReservationDetailRepository(id, communityID, userID, filterByUserID)
	if result.Statue.Error != nil {
		ctx.JSON(http.StatusNotFound, model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "找不到該預約"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result.Result})
}

// AdminDeleteReservation 管理員強制取消/軟刪除預約
// @Summary 管理員強制取消預約
// @Description 管理員軟刪除同社區預約，住戶列表將不再顯示該預約。
// @Tags Reservation
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "預約 ID"
// @Success 200 {object} model.RequestMessage "預約刪除成功"
// @Failure 400 {object} model.Response400Error "無效的預約 ID"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 403 {object} model.Response403Error "權限不足"
// @Failure 404 {object} model.Response404Error "找不到預約"
// @Router /api/v1/admin/reservations/{id}/delete [patch]
func (c *ReservationController) AdminDeleteReservation(ctx *gin.Context) {
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
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid reservation ID"))
		return
	}

	if err := repository.DeleteReservationRepository(id, communityID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "找不到該預約"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "刪除預約失敗"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "預約刪除成功"})
}

// CreateReservation 建立設施預約
// @Summary 建立設施預約
// @Description 住戶針對特定社區設施進行預約。系統會根據設施規則（如黑名單、預約時段重疊等）進行驗證。
// @Tags Reservation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param community_id header string true "社區 ID"
// @Param facility_id path uint64 true "設施 ID"
// @Param reservation body reservationModel.CreateReservationRequest true "預約詳細資料"
// @Success 201 {object} model.RequestMessage "預約建立成功"
// @Failure 400 {object} model.Response400Error "無效的設施 ID 或輸入資料"
// @Failure 403 {object} model.Response403Error "設施不開放或權限不足"
// @Failure 404 {object} model.Response404Error "設施不存在"
// @Failure 409 {object} model.Response409Error "時段衝突或違反設施規則"
// @Router /api/v1/facilities/{facility_id}/reservations [post]
func (c *ReservationController) CreateReservation(ctx *gin.Context) {
	communityID, _ := ctx.Get("community_id")
	userID, _ := ctx.Get("user_id")

	facilityIDStr := ctx.Param("facility_id")
	facilityID, err := strconv.ParseUint(facilityIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid facility ID"))
		return
	}

	var req reservationModel.CreateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid input"))
		return
	}

	// 取得 facility info & rule (檢查是否存在並同屬此社區)
	var facility facilitydb.FacilityInfo
	if err := database.DB.Where("id = ? AND community_id = ?", facilityID, communityID).First(&facility).Error; err != nil {
		ctx.JSON(http.StatusNotFound, model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "設施不存在或權限不足"))
		return
	}
	var rule facilitydb.FacilityRule
	database.DB.Where("facility_id = ?", facilityID).First(&rule) // 若無自帶一筆空的也可以

	if !rule.IsActive || facility.Status != "active" {
		ctx.JSON(http.StatusForbidden, model.NewErrorResponse(ctx, http.StatusForbidden, utilsErr.ErrForbidden, "該設施目前不開放預約"))
		return
	}

	res := &facilitydb.FacilityReservation{
		FacilityID:      facilityID,
		UserID:          userID.(string),
		CommunityID:     communityID.(uint64),
		HomeID:          req.HomeID,
		ReservationDate: req.ReservationDate,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		PeopleCount:     req.PeopleCount,
		Remark:          req.Remark,
	}

	result := repository.CreateReservationRepository(res, &rule)
	if result.Statue.Error != nil {
		ctx.JSON(http.StatusConflict, model.NewErrorResponse(ctx, http.StatusConflict, utilsErr.ErrConflict, result.Statue.Error.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "預約建立成功",
		"data":    result.Result,
	})
}

// CancelReservation 取消預約
// @Summary 取消預約
// @Description 住戶主動取消已建立的設施預約。需提供取消原因。
// @Tags Reservation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "預約紀錄 ID"
// @Param cancel body reservationModel.CancelReservationRequest true "取消原因"
// @Success 200 {object} model.RequestMessage "預約已成功取消"
// @Failure 400 {object} model.Response400Error "無效的 ID 或取消失敗"
// @Router /api/v1/reservations/{id}/cancel [patch]
func (c *ReservationController) CancelReservation(ctx *gin.Context) {
	resIDStr := ctx.Param("id")
	resID, err := strconv.ParseUint(resIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid reservation ID"))
		return
	}
	userID, _ := ctx.Get("user_id")

	var req reservationModel.CancelReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid input"))
		return
	}

	err = repository.CancelReservationRepository(resID, userID.(string), req.CancelReason)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "預約已成功取消"})
}

// Reschedule 申請預約改期
// @Summary 申請預約改期
// @Description 住戶針對現有的預約申請變更日期或時段。申請後需經由管理員審核。
// @Tags Reservation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "預約紀錄 ID"
// @Param reschedule body reservationModel.RescheduleRequest true "新的預約時段與原因"
// @Success 201 {object} model.RequestMessage "改期申請已送出"
// @Failure 400 {object} model.Response400Error "無效的 ID 或輸入資料"
// @Failure 404 {object} model.Response404Error "找不到預約紀錄"
// @Failure 409 {object} model.Response409Error "改期衝突或處理失敗"
// @Router /api/v1/reservations/{id}/reschedule [post]
func (c *ReservationController) Reschedule(ctx *gin.Context) {
	resIDStr := ctx.Param("id")
	resID, err := strconv.ParseUint(resIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid ID"))
		return
	}
	userID, _ := ctx.Get("user_id")

	var req reservationModel.RescheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid input"))
		return
	}

	// 從 URL 得知要針對哪一筆預約進行改期。我們需自己撈出 reservation.FacilityID 以便傳給 Repo
	var res facilitydb.FacilityReservation
	if err := database.DB.Where("id = ? AND user_id = ?", resID, userID).First(&res).Error; err != nil {
		ctx.JSON(http.StatusNotFound, model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "找不到您的預約"))
		return
	}

	rescheduleRecord := &facilitydb.RescheduleRequest{
		ReservationID: resID,
		NewDate:       req.NewDate,
		NewStartTime:  req.NewStartTime,
		NewEndTime:    req.NewEndTime,
		Reason:        req.Reason,
	}

	err = repository.CreateRescheduleRequestRepository(rescheduleRecord, userID.(string), res.FacilityID)
	if err != nil {
		ctx.JSON(http.StatusConflict, model.NewErrorResponse(ctx, http.StatusConflict, utilsErr.ErrConflict, "改期失敗: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "改期申請已送出，等待審核"})
}

// AdminApproveReschedule 管理員審核改期申請
// @Summary 審核預約改期 (管理員)
// @Description 管理員針對住戶提出的改期申請進行核准或拒絕。
// @Tags Reservation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "改期申請紀錄 ID"
// @Param approve body reservationModel.ApproveRescheduleRequest true "審核結果與意見"
// @Success 200 {object} model.RequestMessage "改期申請處理完成"
// @Failure 400 {object} model.Response400Error "無效的 ID 或輸入資料"
// @Failure 500 {object} model.Response500Error "審核處理失敗"
// @Router /api/v1/admin/reschedule/{id}/approve [patch]
func (c *ReservationController) AdminApproveReschedule(ctx *gin.Context) {
	reqIDStr := ctx.Param("id")
	reqID, err := strconv.ParseUint(reqIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid ID"))
		return
	}

	var req reservationModel.ApproveRescheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "Invalid input"))
		return
	}

	err = repository.ApproveRescheduleRequestRepository(reqID, req.IsApprove, req.AdminComment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "審核處理失敗: "+err.Error()))
		return
	}

	status := "已拒絕"
	if req.IsApprove {
		status = "已核准"
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "改期申請" + status})
}
