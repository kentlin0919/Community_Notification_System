package community

// CommunityListQuery 描述社區列表的查詢條件
type CommunityListQuery struct {
	Municipality string `form:"municipality" example:"新北市"`
	District     string `form:"district" example:"淡水區"`
	PostalCode   *int   `form:"postal_code" example:"251"`
	Keyword      string `form:"keyword" example:"甜水郡"`
	Page         int    `form:"page" example:"1"`
	PageSize     int    `form:"page_size" example:"20"`
}

// CommunitySummary 描述單一社區的基本資訊
type CommunitySummary struct {
	CommunityID   uint64 `json:"community_id" example:"1"`
	PostalCode    int    `json:"postal_code" example:"251"`
	Municipality  string `json:"municipality" example:"新北市"`
	District      string `json:"district" example:"淡水區"`
	RoadName      string `json:"road_name" example:"濱海路一段"`
	LaneNumber    int    `json:"lane_number" example:"306"`
	AlleyNumber   int    `json:"alley_number" example:"0"`
	CommunityName string `json:"community_name" example:"甜水郡社區"`
	Address       string `json:"address" example:"251新北市淡水區濱海路一段306巷"`
}

// CommunityListResponse 用於回傳社區清單與統計
type CommunityListResponse struct {
	Total       int64              `json:"total" example:"1"`
	Communities []CommunitySummary `json:"communities"`
}

type CommunityRegister struct {
	PermissionID  string `gorm:"uniqueIndex;not null" json:"permission_id"`
	PostalCode    int    `json:"postal_code" gorm:"type:int;not null;comment:郵遞區號"`
	Municipality  string `json:"municipality" gorm:"type:varchar(10);not null;comment:縣市"`
	District      string `json:"district" gorm:"type:varchar(10);not null;comment:鄉鎮市區"`
	RoadName      string `json:"road_name" gorm:"type:varchar(50);not null;comment:路名"`
	LaneNumber    int    `json:"lane_number" gorm:"type:int;not null;comment:巷弄號碼"`
	AlleyNumber   int    `json:"alley_number" gorm:"type:int;not null;comment:巷弄號碼"`
	CommunityName string `json:"community_name" gorm:"type:varchar(50);not null;comment:社區名稱"`
	Address       string `json:"address" gorm:"type:varchar(100);not null;comment:地址"`
}
