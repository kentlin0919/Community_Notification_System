package communityManager

// / 主要處理社區管理相關的請求
// / 基礎項目包含社區ID、郵遞區號、鄉鎮市區、路名、巷弄號碼、社區名稱及地址等
type CommunityManagerController struct{}

// NewCommunityTableController 建構函式，建立一個新的 CommunityTableController
func NewCommunityTableController() *CommunityManagerController {
	return &CommunityManagerController{}
}
