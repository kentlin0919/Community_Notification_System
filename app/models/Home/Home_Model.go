package home

// HomeResponse 首頁聚合回應結構
type HomeResponse struct {
	Role string      `json:"role" example:"resident"`
	Data interface{} `json:"data"`
}

// ── 住戶 (PermissionID 8-9) ──

// ResidentCounts 住戶首頁摘要數字
type ResidentCounts struct {
	UnreadMessages       int64 `json:"unread_messages"`
	UnclaimedParcels     int64 `json:"unclaimed_parcels"`
	UpcomingReservations int64 `json:"upcoming_reservations"`
}

// MessageSummary 訊息摘要（首頁用）
type MessageSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Category   string `json:"category"`
	IsRead     bool   `json:"is_read"`
	CreateTime string `json:"create_time"`
}

// ReservationSummary 預約摘要（首頁用）
type ReservationSummary struct {
	ID              uint64 `json:"id"`
	FacilityName    string `json:"facility_name"`
	ReservationDate string `json:"date"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	Status          string `json:"status"`
}

// ResidentDashboard 住戶首頁資料
type ResidentDashboard struct {
	Counts              ResidentCounts       `json:"counts"`
	RecentMessages      []MessageSummary     `json:"recent_messages"`
	UpcomingReservations []ReservationSummary `json:"upcoming_reservations"`
}

// ── 社區管理員 (PermissionID 3-7) ──

// AdminStats 管理員首頁統計數字
type AdminStats struct {
	ResidentCount     int64 `json:"resident_count"`
	TodayReservations int64 `json:"today_reservations"`
	UnhandledParcels  int64 `json:"unhandled_parcels"`
}

// AdminDashboard 社區管理員首頁資料
type AdminDashboard struct {
	Stats             AdminStats           `json:"stats"`
	TodayReservations []ReservationSummary `json:"today_reservations"`
	RecentMessages    []MessageSummary     `json:"recent_messages"`
}

// ── 超級管理員 (PermissionID 1-2) ──

// PlatformStats 平台統計數字
type PlatformStats struct {
	CommunityCount     int64 `json:"community_count"`
	UserCount          int64 `json:"user_count"`
	PendingCommunities int64 `json:"pending_communities"`
}

// PendingCommunitySummary 待審核社區摘要
type PendingCommunitySummary struct {
	ID            uint64 `json:"id"`
	CommunityName string `json:"community_name"`
	ApplicantName string `json:"applicant_name"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// SuperAdminDashboard 超級管理員首頁資料
type SuperAdminDashboard struct {
	PlatformStats      PlatformStats             `json:"platform_stats"`
	PendingCommunities []PendingCommunitySummary  `json:"pending_communities"`
}
