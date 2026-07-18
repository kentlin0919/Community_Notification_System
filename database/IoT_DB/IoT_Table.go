package iot_db

import (
	"gorm.io/gorm"
	"log"
)

type IoTDBController struct{}

func NewIoTDBController() *IoTDBController {
	return &IoTDBController{}
}

func (i *IoTDBController) IoTTable(db *gorm.DB) {
	err := db.AutoMigrate(
		&IoTDevice{},
		&VisitorPass{},
		&AccessLog{},
	)
	if err != nil {
		log.Fatalf("IoT 模組資料表建立失敗: %v", err)
	}
	log.Println("IoT 模組資料表建立/遷移成功")
}
