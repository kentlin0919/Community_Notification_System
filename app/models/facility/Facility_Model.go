package facility

type CreateFacilityRequest struct {
	CommunityID  uint64 `json:"community_id" binding:"required" example:"1"`
	Name         string `json:"name" binding:"required,max=100" example:"健身房"`
	FacilityType string `json:"facility_type" binding:"required,max=50" example:"運動設施"`
	Location     string `json:"location" example:"B1 棟"`
	Description  string `json:"description" example:"提供跑步機、啞鈴等健身設備。"`
	Status       string `json:"status" binding:"omitempty,oneof=draft active inactive maintenance" example:"active"`
	CoverImage   string `json:"cover_image" example:"https://example.com/facility.jpg"`
}
