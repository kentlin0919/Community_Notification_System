package maintenance_db

import (
	"Community_Notification_System/pkg/databasepkg"

	"gorm.io/gorm"
)

type MaintenanceController struct{}

func NewMaintenanceController() *MaintenanceController {
	return &MaintenanceController{}
}

func (m *MaintenanceController) MaintenanceTable(DB *gorm.DB) {
	databasepkg.NewCreateTableController().Base_Create_Table(DB, &MaintenanceReports{}, "maintenance_reports")
}
