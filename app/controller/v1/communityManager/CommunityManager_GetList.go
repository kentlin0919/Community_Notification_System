package communityManager

import (
	"net/http"

	communityModel "Community_Notification_System/app/models/community"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/community"

	"github.com/gin-gonic/gin"
)

// CommunityManager_GetList 取得社區基本資料清單
// @Summary 取得社區列表
// @Description 依縣市、行政區、郵遞區號或關鍵字等條件查詢社區基本資料清單，並支援分頁查詢。
// @Tags CommunityManager
// @Accept json
// @Produce json
// @Param municipality query string false "縣市過濾（例如：新北市）"
// @Param district query string false "行政區過濾（例如：淡水區）"
// @Param postal_code query int false "郵遞區號過濾（例如：251）"
// @Param keyword query string false "關鍵字搜尋（社區名稱或地址）"
// @Param page query int false "頁碼（從 1 開始，預設 1）"
// @Param page_size query int false "每頁筆數（預設 20，最大 100）"
// @Success 200 {object} communityModel.CommunityListResponse "成功取得社區清單"
// @Failure 400 {object} model.ErrorRequest "請求參數錯誤"
// @Failure 401 {object} model.ErrorRequest "未授權"
// @Failure 500 {object} model.ErrorRequest "系統錯誤"
// @Security BearerAuth
// @Router /api/v1/community/getlist [get]
func (c *CommunityManagerController) CommunityManager_GetList(ctx *gin.Context) {
	var query communityModel.CommunityListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "請求參數錯誤")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	} else if query.PageSize > 100 {
		query.PageSize = 100
	}

	repoResult := repository.CommunityListRepository(&query)
	if repoResult.Statue.Error != nil {
		errorModel := model.NewErrorRequest(http.StatusInternalServerError, "取得社區資料失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	response := communityModel.CommunityListResponse{
		Total: repoResult.Result.Total,
	}
	response.Communities = make([]communityModel.CommunitySummary, len(repoResult.Result.Items))
	for idx, item := range repoResult.Result.Items {
		response.Communities[idx] = communityModel.CommunitySummary{
			CommunityID:   item.Community_id,
			PostalCode:    item.PostalCode,
			Municipality:  item.Municipality,
			District:      item.District,
			RoadName:      item.RoadName,
			LaneNumber:    item.LaneNumber,
			AlleyNumber:   item.AlleyNumber,
			CommunityName: item.Community_name,
			Address:       item.Address,
		}
	}

	ctx.JSON(http.StatusOK, response)
}
