package platform

import (
	"net/http"

	"Community_Notification_System/app/models/model"
	platformModel "Community_Notification_System/app/models/platform"
	repository "Community_Notification_System/app/repositories/platform"
	"Community_Notification_System/utils"

	utilsErr "Community_Notification_System/utils/errors"

	"github.com/gin-gonic/gin"
)

// SystemConfig 檢查 App 版本並回傳是否需要強制更新
// @Summary 檢查 App 版本
// @Description 依據 os 與目前版本，回傳最低版本、最新版本與是否需強制更新
// @Tags Platform
// @Accept json
// @Produce json
// @Param os query string true "作業系統 (ios/android)"
// @Param version query string true "目前 App 版本"
// @Success 200 {object} platformModel.SystemConfigResponse "版本檢查結果"
// @Failure 404 {object} model.Response404Error "查無該平台版本設定"
// @Router /api/v1/system/config [get]
func (p *PlatformController) SystemConfig(ctx *gin.Context) {
	os := ctx.Query("os")
	version := ctx.Query("version")

	result := repository.GetPlatformVersionConfig(os)
	if result.Statue.Error != nil || result.Result == nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrNotFound, "查無該平台版本設定")
		ctx.JSON(http.StatusNotFound, errorModel)
		return
	}

	forceUpdate := utils.CompareVersions(version, result.Result.MinVersion) < 0

	ctx.JSON(http.StatusOK, platformModel.SystemConfigResponse{
		MinVersion:    result.Result.MinVersion,
		LatestVersion: result.Result.LatestVersion,
		ForceUpdate:   forceUpdate,
	})
}
