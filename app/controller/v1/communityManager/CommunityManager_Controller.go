package communityManager

import (
	repository "Community_Notification_System/app/repositories/community"
	communitydb "Community_Notification_System/database/Community_DB"
	"errors"

	"gorm.io/gorm"
)

// / 主要處理社區管理相關的請求
// / 基礎項目包含社區ID、郵遞區號、鄉鎮市區、路名、巷弄號碼、社區名稱及地址等
type CommunityManagerController struct{}

// NewCommunityTableController 建構函式，建立一個新的 CommunityTableController
func NewCommunityTableController() *CommunityManagerController {
	return &CommunityManagerController{}
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
