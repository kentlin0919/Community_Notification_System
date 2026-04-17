package communitydb

import (
	"time"
)

type CommunityRegisterApplication struct {
	ID                uint64    `json:"id" gorm:"primarykey;autoIncrement"`
	Status            string    `json:"status" gorm:"type:varchar(20);default:'pending';not null"` // pending, approved, rejected
	PostalCode        int       `json:"postal_code"`
	Municipality      string    `json:"municipality" gorm:"type:varchar(20);not null"`
	District          string    `json:"district" gorm:"type:varchar(20);not null"`
	RoadName          string    `json:"road_name" gorm:"type:varchar(50);not null"`
	LaneNumber        int       `json:"lane_number"`
	AlleyNumber       int       `json:"alley_number"`
	CommunityName     string    `json:"community_name" gorm:"type:varchar(100);not null"`
	Address           string    `json:"address" gorm:"type:varchar(255);not null"`
	AdminName         string    `json:"admin_name" gorm:"type:varchar(50);not null"`
	AdminEmail        string    `json:"admin_email" gorm:"type:varchar(100);not null"`
	AdminPasswordHash string    `json:"-" gorm:"type:varchar(255);not null"` // 保存密碼 Hash
	AdminPhone        string    `json:"admin_phone" gorm:"type:varchar(20)"`
	ApplicantName     string    `json:"applicant_name" gorm:"type:varchar(50);not null"`
	ApplicantEmail    string    `json:"applicant_email" gorm:"type:varchar(100);not null"`
	ApplicantPhone    string    `json:"applicant_phone" gorm:"type:varchar(20)"`
	Remark            string    `json:"remark" gorm:"type:text"`
	ReviewedBy        *string   `json:"reviewed_by"` // Super admin's User ID
	ReviewedAt        *time.Time `json:"reviewed_at"`
	RejectReason      string    `json:"reject_reason" gorm:"type:text"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
