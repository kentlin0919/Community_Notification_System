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

// SystemConfigResponse 定義版本檢查回應格式
type SystemConfigResponse struct {
	MinVersion    string `json:"min_version"`
	LatestVersion string `json:"latest_version"`
	ForceUpdate   bool   `json:"force_update"`
}
