package visitor_db

import (
	"Community_Notification_System/pkg/databasepkg"

	"gorm.io/gorm"
)

type VisitorController struct{}

func NewVisitorController() *VisitorController {
	return &VisitorController{}
}

func (v *VisitorController) VisitorTable(DB *gorm.DB) {
	databasepkg.NewCreateTableController().Base_Create_Table(DB, &VisitorInfo{}, "visitor_info")
}
