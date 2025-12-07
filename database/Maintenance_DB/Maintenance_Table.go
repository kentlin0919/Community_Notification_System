package maintenance_db

import (
	"Community_Notification_System/pkg/databasePkg"

	"gorm.io/gorm"
)

type MaintenanceController struct{}

func NewMaintenanceController() *MaintenanceController {
	return &MaintenanceController{}
}

func (m *MaintenanceController) MaintenanceTable(DB *gorm.DB) {
	databasePkg.NewCreateTableController().Base_Create_Table(DB, &MaintenanceReports{}, "maintenance_reports")
}
