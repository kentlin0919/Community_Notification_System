package visitor_db

import (
	"Community_Notification_System/pkg/databasePkg"

	"gorm.io/gorm"
)

type VisitorController struct{}

func NewVisitorController() *VisitorController {
	return &VisitorController{}
}

func (v *VisitorController) VisitorTable(DB *gorm.DB) {
	databasePkg.NewCreateTableController().Base_Create_Table(DB, &VisitorInfo{}, "visitor_info")
}
