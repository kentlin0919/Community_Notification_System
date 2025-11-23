package platform

import (
	"net/http"

	"Community_Notification_System/app/models/model"
	platformModel "Community_Notification_System/app/models/platform"
	repository "Community_Notification_System/app/repositories/platform"

	"github.com/gin-gonic/gin"
)

// Platform_GetList 取得平台列表
// @Summary 取得平台列表
// @Description 取得所有可用平台名稱，供前端顯示於選項列表。
// @Tags Platform
// @Accept json
// @Produce json
// @Success 200 {object} platformModel.PlatformListResponse "成功取得平台清單"
// @Failure 500 {object} model.ErrorRequest "取得平台資料失敗"
// @Router /api/v1/platform/getlist [get]
func (p *PlatformController) Platform_GetList(ctx *gin.Context) {
	repoResult := repository.PlatformRepository()

	if repoResult.Statue.Error != nil {
		errorModel := model.NewErrorRequest(http.StatusInternalServerError, "取得平台資料失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	response := platformModel.PlatformListResponse{
		Total:     repoResult.Result.Total,
		Platforms: make([]platformModel.PlatformSummary, len(repoResult.Result.Items)),
	}

	for idx, item := range repoResult.Result.Items {
		response.Platforms[idx] = platformModel.PlatformSummary{Name: item}
	}

	ctx.JSON(http.StatusOK, response)
}
