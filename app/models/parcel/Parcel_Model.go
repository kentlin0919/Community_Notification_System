package parcel

import "time"

// CreateParcelRequest 管理員代收登錄請求
type CreateParcelRequest struct {
	HomeID         uint64 `json:"home_id" binding:"required" example:"1"`
	CourierCompany string `json:"courier_company" binding:"required" example:"黑貓宅急便"`
	TrackingNumber string `json:"tracking_number" binding:"required" example:"1234567890"`
	Remark         string `json:"remark" example:"貴重物品"`
}

// ParcelListResponse 包裹列表回應
type ParcelListResponse struct {
	ID             uint64     `json:"id" example:"1"`
	CommunityID    uint64     `json:"community_id" example:"1"`
	HomeID         uint64     `json:"home_id" example:"1"`
	CourierCompany string     `json:"courier_company" example:"黑貓宅急便"`
	TrackingNumber string     `json:"tracking_number" example:"1234567890"`
	Status         int        `json:"status" example:"0"` // 0: 待領取, 1: 已領取
	Remark         string     `json:"remark" example:"貴重物品"`
	ReceivedAt     time.Time  `json:"received_at" example:"2025-05-10T08:00:00Z"`
	PickedUpAt     *time.Time `json:"picked_up_at" example:"2025-05-10T10:00:00Z"`
	CreatedBy      string     `json:"created_by" example:"admin-uuid"`
}

// ParcelPickupResponse 領取成功回應
type ParcelPickupResponse struct {
	Message    string    `json:"message" example:"包裹領取成功"`
	PickedUpAt time.Time `json:"picked_up_at" example:"2025-05-10T10:00:00Z"`
}
