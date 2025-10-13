package platform

// PlatformSummary 描述單一平台資訊
type PlatformSummary struct {
	Name string `json:"name"`
}

// PlatformListResponse 定義平台清單回應格式
type PlatformListResponse struct {
	Total     int64             `json:"total"`
	Platforms []PlatformSummary `json:"platforms"`
}
