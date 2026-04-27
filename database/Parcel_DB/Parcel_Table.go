package parcel_db

import (
	"Community_Notification_System/pkg/common"

	"gorm.io/gorm"
)

type ParcelDBController struct{}

func NewParcelDBController() *ParcelDBController {
	return &ParcelDBController{}
}

func (p *ParcelDBController) ParcelTable(DB *gorm.DB) {
	common.NewCreateTableController().Base_Create_Table(DB, &ParcelInfo{}, "parcel_info")
}
