package community

// RegisterApplicationRequest 描述社區申請的 JSON Payload
type RegisterApplicationRequest struct {
	PostalCode     int    `json:"postal_code" example:"251" binding:"required"`
	Municipality   string `json:"municipality" example:"新北市" binding:"required"`
	District       string `json:"district" example:"淡水區" binding:"required"`
	RoadName       string `json:"road_name" example:"濱海路一段" binding:"required"`
	LaneNumber     int    `json:"lane_number" example:"306"`
	AlleyNumber    int    `json:"alley_number" example:"0"`
	CommunityName  string `json:"community_name" example:"甜水郡社區" binding:"required"`
	Address        string `json:"address" example:"251新北市淡水區濱海路一段306巷" binding:"required"`
	AdminName      string `json:"admin_name" example:"AdminName" binding:"required"`
	AdminEmail     string `json:"admin_email" example:"admin@example.com" binding:"required,email"`
	AdminPassword  string `json:"admin_password" example:"password123" binding:"required,min=6"`
	AdminPhone     string `json:"admin_phone" example:"0912345678"`
	ApplicantName  string `json:"applicant_name" example:"ApplicantName" binding:"required"`
	ApplicantEmail string `json:"applicant_email" example:"applicant@example.com" binding:"required,email"`
	ApplicantPhone string `json:"applicant_phone" example:"0987654321"`
	Remark         string `json:"remark" example:"申請建立測試社區"`
}

type ApproveApplicationRequest struct {
	// 核可通常不需要額外 Payload，若是需要補充備註可留此結構，或者直接使用 url param
}

type RejectApplicationRequest struct {
	RejectReason string `json:"reject_reason" example:"資訊不全，請重新申請"`
}
