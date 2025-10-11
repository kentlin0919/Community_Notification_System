package communityManager

import (
	communityModel "Community_Notification_System/app/models/community"
	"Community_Notification_System/app/models/model"

	repository "Community_Notification_System/app/repositories/community"
	communitydb "Community_Notification_System/database/Community_DB"

	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommunityManager_Register 新增社區基本資料清單
// @Summary 新增社區
// @Description 新增社區基本資料，欄位說明如下：
// @Description - permission_id：綁定社區管理權限的唯一識別字串。
// @Description - postal_code：社區所在地的郵遞區號，以數字表示便於郵遞與查詢。
// @Description - municipality：社區所在的縣市名稱。
// @Description - district：社區所在的鄉鎮市區名稱。
// @Description - road_name：社區主要道路名稱，例如路、街或段。
// @Description - lane_number：地址中的巷或弄號碼，若無請填 0。
// @Description - alley_number：更細分的弄號碼或巷內編號，若無請填 0。
// @Description - community_name：社區或大樓的正式名稱。
// @Description - address：完整地址（含門牌號），提供精確位置資訊。
// @Tags CommunityManager
// @Accept json
// @Produce json
// @Param community body communityModel.CommunityRegister true "社區基本資料"
// @Success 200 {object} map[string]string "成功新增社區"
// @Failure 400 {object} model.ErrorRequest "請求參數錯誤"
// @Failure 401 {object} model.ErrorRequest "未授權"
// @Failure 500 {object} model.ErrorRequest "系統錯誤"
// @Security BearerAuth
// @Router /api/v1/community/register [post]
func (c *CommunityManagerController) CommunityManager_Register(ctx *gin.Context) {
	var registerModel communityModel.CommunityRegister

	var community_info = &communitydb.CommunityInfo{}

	// 綁定 JSON 資料，若結構無法對應或缺欄位則直接回傳 400
	if err := ctx.ShouldBindJSON(&registerModel); err != nil {
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "Invalid input")
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	community_info.Municipality = registerModel.Municipality
	community_info.District = registerModel.District
	community_info.PostalCode = registerModel.PostalCode
	community_info.RoadName = registerModel.RoadName
	community_info.LaneNumber = registerModel.LaneNumber
	community_info.AlleyNumber = registerModel.AlleyNumber
	community_info.Community_name = registerModel.CommunityName
	community_info.Address = registerModel.Address

	checkCommunityStatue, err := checkCommunity(community_info)

	if err != nil {
		errorModel := model.NewErrorRequest(http.StatusInternalServerError, "檢查社區資料失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	if !checkCommunityStatue {
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "社區已經存在")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	result := repository.RegisterRepository(community_info)

	if result.Statue.Error != nil {
		errorModel := model.NewErrorRequest(http.StatusInternalServerError, "新增社區失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "社區新增成功",
	})

}

func checkCommunity(communityData *communitydb.CommunityInfo) (bool, error) {
	result := repository.CommunityOneRepository(*communityData)
	if result.Statue.Error != nil {
		if errors.Is(result.Statue.Error, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, result.Statue.Error
	}

	return result.Statue.RowsAffected == 0, nil
}
