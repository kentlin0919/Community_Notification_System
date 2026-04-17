package reservation

type CreateReservationRequest struct {
	HomeID          uint64 `json:"home_id"` // 可選，看業務邏輯
	ReservationDate string `json:"reservation_date" binding:"required"` // YYYY-MM-DD
	StartTime       string `json:"start_time" binding:"required"`       // HH:mm
	EndTime         string `json:"end_time" binding:"required"`         // HH:mm
	PeopleCount     int    `json:"people_count" binding:"required,min=1"`
	Remark          string `json:"remark"`
}

type CancelReservationRequest struct {
	CancelReason string `json:"cancel_reason"`
}

type RescheduleRequest struct {
	NewDate      string `json:"new_date" binding:"required"`
	NewStartTime string `json:"new_start_time" binding:"required"`
	NewEndTime   string `json:"new_end_time" binding:"required"`
	Reason       string `json:"reason"`
}

type ApproveRescheduleRequest struct {
	AdminComment string `json:"admin_comment"`
	IsApprove    bool   `json:"is_approve"` // true: 核准, false: 拒絕
}
