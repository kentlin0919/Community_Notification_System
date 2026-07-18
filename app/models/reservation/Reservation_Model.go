package reservation

type CreateReservationRequest struct {
	HomeID          uint64 `json:"home_id" example:"1"`                                      // 可選，看業務邏輯
	ReservationDate string `json:"reservation_date" binding:"required" example:"2025-05-20"` // YYYY-MM-DD
	StartTime       string `json:"start_time" binding:"required" example:"14:00"`            // HH:mm
	EndTime         string `json:"end_time" binding:"required" example:"15:00"`              // HH:mm
	PeopleCount     int    `json:"people_count" binding:"required,min=1" example:"2"`
	Remark          string `json:"remark" example:"使用跑步機"`
}

type CancelReservationRequest struct {
	CancelReason string `json:"cancel_reason" example:"臨時有事無法前往"`
}

type RescheduleRequest struct {
	NewDate      string `json:"new_date" binding:"required" example:"2025-05-21"`
	NewStartTime string `json:"new_start_time" binding:"required" example:"10:00"`
	NewEndTime   string `json:"new_end_time" binding:"required" example:"11:00"`
	Reason       string `json:"reason" example:"原本時段突然要開會"`
}

type ApproveRescheduleRequest struct {
	AdminComment string `json:"admin_comment" example:"已核准改期申請"`
	IsApprove    bool   `json:"is_approve" example:"true"` // true: 核准, false: 拒絕
}
