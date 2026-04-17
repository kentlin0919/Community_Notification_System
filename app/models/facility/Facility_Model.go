package facility

type CreateFacilityRequest struct {
	CommunityID  uint64 `json:"community_id" binding:"required"`
	Name         string `json:"name" binding:"required,max=100"`
	FacilityType string `json:"facility_type" binding:"required,max=50"`
	Location     string `json:"location"`
	Description  string `json:"description"`
	Status       string `json:"status" binding:"omitempty,oneof=draft active inactive maintenance"`
	CoverImage   string `json:"cover_image"`
}
