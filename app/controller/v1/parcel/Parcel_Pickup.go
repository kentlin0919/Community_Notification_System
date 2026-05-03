package parcel

import (
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/parcel"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PickupParcel 標記包裹為已領取
// @Summary 標記包裹已領取
// @Description 管理員確認住戶領取包裹後，將包裹狀態更新為「已領取」
// @Tags Parcel Management
// @Accept json
// @Produce json
// @Param id path int true "包裹 ID"
// @Success 200 {object} map[string]interface{} "包裹領取成功"
// @Failure 400 {object} model.Response400Error "請求參數錯誤或重複領取"
// @Failure 401 {object} model.Response401Error "無法取得登入資訊"
// @Failure 404 {object} model.Response404Error "找不到包裹"
// @Failure 500 {object} model.Response500Error "伺服器錯誤"
// @Security BearerAuth
// @Router /api/v1/parcels/{id}/pickup [put]
func (c *ParcelController) PickupParcel(ctx *gin.Context) {
	communityID, exists := ctx.Get("community_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorRequest(http.StatusUnauthorized, "無法存取社區資訊，請重新登入"))
		return
	}

	cID := communityID.(uint64)

	// 解析路徑參數
	parcelIDStr := ctx.Param("id")
	parcelID, err := strconv.ParseUint(parcelIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, "包裹 ID 格式錯誤"))
		return
	}

	res := repository.PickupParcelRepository(cID, parcelID)
	if res.Statue.Error != nil {
		errMsg := res.Statue.Error.Error()
		switch errMsg {
		case "找不到該包裹紀錄":
			ctx.JSON(http.StatusNotFound, model.NewErrorRequest(http.StatusNotFound, errMsg))
		case "該包裹已被領取":
			ctx.JSON(http.StatusBadRequest, model.NewErrorRequest(http.StatusBadRequest, errMsg))
		default:
			ctx.JSON(http.StatusInternalServerError, model.NewErrorRequest(http.StatusInternalServerError, "更新包裹狀態失敗"))
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "包裹領取成功",
		"picked_up_at": res.Result.PickedUpAt,
	})
}
