package parcel

import (
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/parcel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetParcelList 查詢包裹列表
// @Summary 查詢包裹列表
// @Description 查詢指定社區的包裹列表，可透過 home_id 與 status 篩選
// @Tags Parcel Management
// @Accept json
// @Produce json
// @Param home_id query int false "住戶 ID"
// @Param status query int false "狀態 (1=待領取, 2=已領取)"
// @Success 200 {object} map[string]interface{} "查詢成功"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 500 {object} model.Response500Error "伺服器錯誤"
// @Security BearerAuth
// @Router /api/v1/parcels [get]
func (c *ParcelController) GetParcelList(ctx *gin.Context) {
	communityID, exists := ctx.Get("community_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorRequest(http.StatusUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}

	cID := communityID.(uint64)

	// 解析查詢參數
	var homeID uint64
	var status int

	if homeIDStr := ctx.Query("home_id"); homeIDStr != "" {
		parsed, err := strconv.ParseUint(homeIDStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "home_id 格式錯誤"))
			return
		}
		homeID = parsed
	}

	if statusStr := ctx.Query("status"); statusStr != "" {
		parsed, err := strconv.Atoi(statusStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "status 格式錯誤"))
			return
		}
		status = parsed
	}

	res := repository.GetParcelListRepository(cID, homeID, status)
	if res.Statue.Error != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "查詢包裹列表失敗"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "查詢成功",
		"data":    res.Result,
		"total":   len(res.Result),
	})
}
