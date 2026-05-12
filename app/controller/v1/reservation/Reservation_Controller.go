package reservation

import (
	reservationModel "Community_Notification_System/app/models/reservation"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/reservation"
	facilitydb "Community_Notification_System/database/Facility_DB"
	"Community_Notification_System/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReservationController struct{}

func NewReservationController() *ReservationController {
	return &ReservationController{}
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
// @Success 201 {object} models.RequestMessage "預約建立成功"
// @Failure 400 {object} model.Response400Error "無效的設施 ID 或輸入資料"
// @Failure 403 {object} model.Response403Error "設施不開放或權限不足"
// @Failure 404 {object} model.Response404Error "設施不存在"
// @Failure 409 {object} model.Response409Error "時段衝突或違反設施規則"
// @Router /api/v1/reservation/{facility_id} [post]
func (c *ReservationController) CreateReservation(ctx *gin.Context) {
	communityID, _ := ctx.Get("community_id")
	userID, _ := ctx.Get("user_id")

	facilityIDStr := ctx.Param("facility_id")
	facilityID, err := strconv.ParseUint(facilityIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid facility ID"))
		return
	}

	var req reservationModel.CreateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid input"))
		return
	}

	// 取得 facility info & rule (檢查是否存在並同屬此社區)
	var facility facilitydb.FacilityInfo
	if err := database.DB.Where("id = ? AND community_id = ?", facilityID, communityID).First(&facility).Error; err != nil {
		ctx.JSON(http.StatusNotFound, model.NewErrorRequest(http.StatusNotFound, "設施不存在或權限不足"))
		return
	}
	var rule facilitydb.FacilityRule
	database.DB.Where("facility_id = ?", facilityID).First(&rule) // 若無自帶一筆空的也可以

	if !rule.IsActive || facility.Status != "active" {
		ctx.JSON(http.StatusForbidden, model.NewErrorRequest(http.StatusForbidden, "該設施目前不開放預約"))
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
		ctx.JSON(http.StatusConflict, model.NewErrorRequest(http.StatusConflict, result.Statue.Error.Error()))
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
// @Success 200 {object} models.RequestMessage "預約已成功取消"
// @Failure 400 {object} model.Response400Error "無效的 ID 或取消失敗"
// @Router /api/v1/reservation/{id}/cancel [post]
func (c *ReservationController) CancelReservation(ctx *gin.Context) {
	resIDStr := ctx.Param("id")
	resID, err := strconv.ParseUint(resIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid reservation ID"))
		return
	}
	userID, _ := ctx.Get("user_id")

	var req reservationModel.CancelReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid input"))
		return
	}

	err = repository.CancelReservationRepository(resID, userID.(string), req.CancelReason)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, err.Error()))
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
// @Success 201 {object} models.RequestMessage "改期申請已送出"
// @Failure 400 {object} model.Response400Error "無效的 ID 或輸入資料"
// @Failure 404 {object} model.Response404Error "找不到預約紀錄"
// @Failure 409 {object} model.Response409Error "改期衝突或處理失敗"
// @Router /api/v1/reservation/{id}/reschedule [post]
func (c *ReservationController) Reschedule(ctx *gin.Context) {
	resIDStr := ctx.Param("id")
	resID, err := strconv.ParseUint(resIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid ID"))
		return
	}
	userID, _ := ctx.Get("user_id")

	var req reservationModel.RescheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid input"))
		return
	}

	// 從 URL 得知要針對哪一筆預約進行改期。我們需自己撈出 reservation.FacilityID 以便傳給 Repo
	var res facilitydb.FacilityReservation
	if err := database.DB.Where("id = ? AND user_id = ?", resID, userID).First(&res).Error; err != nil {
		ctx.JSON(http.StatusNotFound, model.NewErrorRequest(http.StatusNotFound, "找不到您的預約"))
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
		ctx.JSON(http.StatusConflict, model.NewErrorRequest(http.StatusConflict, "改期失敗: "+err.Error()))
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
// @Success 200 {object} models.RequestMessage "改期申請處理完成"
// @Failure 400 {object} model.Response400Error "無效的 ID 或輸入資料"
// @Failure 500 {object} model.Response500Error "審核處理失敗"
// @Router /api/v1/reservation/admin/reschedule/{id} [post]
func (c *ReservationController) AdminApproveReschedule(ctx *gin.Context) {
	reqIDStr := ctx.Param("id")
	reqID, err := strconv.ParseUint(reqIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid ID"))
		return
	}

	var req reservationModel.ApproveRescheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "Invalid input"))
		return
	}

	err = repository.ApproveRescheduleRequestRepository(reqID, req.IsApprove, req.AdminComment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "審核處理失敗: "+err.Error()))
		return
	}

	status := "已拒絕"
	if req.IsApprove {
		status = "已核准"
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "改期申請" + status})
}
