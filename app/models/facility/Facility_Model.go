package facility

import (
	facility_db "Community_Notification_System/database/Facility_DB"
)

type CreateFacilityRequest struct {
	CommunityID  uint64 `json:"community_id" binding:"required" example:"1"`
	Name         string `json:"name" binding:"required,max=100" example:"健身房"`
	FacilityType string `json:"facility_type" binding:"required,max=50" example:"運動設施"`
	Location     string `json:"location" example:"B1 棟"`
	Description  string `json:"description" example:"提供跑步機、啞鈴等健身設備。"`
	Status       string `json:"status" binding:"omitempty,oneof=draft active inactive maintenance" example:"active"`
	CoverImage   string `json:"cover_image" example:"https://example.com/facility.jpg"`
}

type UpdateFacilityRequest struct {
	Name                string `json:"name" binding:"required,max=100" example:"健身房"`
	FacilityType        string `json:"facility_type" binding:"required,max=50" example:"運動設施"`
	Location            string `json:"location" example:"B1 棟"`
	Description         string `json:"description" example:"提供跑步機、啞鈴等健身設備。"`
	Status              string `json:"status" binding:"omitempty,oneof=draft active inactive maintenance" example:"active"`
	CoverImage          string `json:"cover_image" example:"https://example.com/facility.jpg"`
	SlotMinutes         *int   `json:"slot_minutes" binding:"omitempty,min=1" example:"60"`
	MaxAdvanceDays      *int   `json:"max_advance_days" binding:"omitempty,min=0" example:"7"`
	CancelBeforeHours   *int   `json:"cancel_before_hours" binding:"omitempty,min=0" example:"24"`
	AutoApprove         *bool  `json:"auto_approve" example:"false"`
	ReservationRuleOpen *bool  `json:"reservation_rule_open" example:"true"`
}

// FacilityListResponse 設施列表回傳格式
type FacilityListResponse struct {
	Data []facility_db.FacilityInfo `json:"data"`
}

// FacilityDetailResponse 設施詳情回傳格式
type FacilityDetailResponse struct {
	Data facility_db.FacilityInfo `json:"data"`
}

