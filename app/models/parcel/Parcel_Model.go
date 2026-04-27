package parcel

import "time"

// CreateParcelRequest 管理員代收登錄請求
type CreateParcelRequest struct {
	HomeID         uint64 `json:"home_id" binding:"required"`
	CourierCompany string `json:"courier_company" binding:"required"`
	TrackingNumber string `json:"tracking_number" binding:"required"`
	Remark         string `json:"remark"`
}

// ParcelListResponse 包裹列表回應
type ParcelListResponse struct {
	ID             uint64     `json:"id"`
	CommunityID    uint64     `json:"community_id"`
	HomeID         uint64     `json:"home_id"`
	CourierCompany string     `json:"courier_company"`
	TrackingNumber string     `json:"tracking_number"`
	Status         int        `json:"status"`
	Remark         string     `json:"remark"`
	ReceivedAt     time.Time  `json:"received_at"`
	PickedUpAt     *time.Time `json:"picked_up_at"`
	CreatedBy      string     `json:"created_by"`
}

// ParcelPickupResponse 領取成功回應
type ParcelPickupResponse struct {
	Message    string    `json:"message"`
	PickedUpAt time.Time `json:"picked_up_at"`
}
