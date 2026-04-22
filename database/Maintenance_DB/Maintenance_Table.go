package maintenance_db

import (
	"log"
	"gorm.io/gorm"
)

type MaintenanceDBController struct{}

func NewMaintenanceDBController() *MaintenanceDBController {
	return &MaintenanceDBController{}
}

func (m *MaintenanceDBController) MaintenanceTable(db *gorm.DB) {
	err := db.AutoMigrate(
		&AssetTag{},
		&MaintenanceTicket{},
	)
	if err != nil {
		log.Fatalf("維修模組資料表建立失敗: %v", err)
	}
	log.Println("維修模組資料表建立/遷移成功")
}
