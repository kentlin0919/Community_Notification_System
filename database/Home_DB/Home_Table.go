package Home_db

import (
	"Community_Notification_System/pkg/common"
	"errors"
	"log"

	"gorm.io/gorm"
)

type UserHomeTablesController struct{}

func NewUserHomeTableController() *UserHomeTablesController {
	return &UserHomeTablesController{}
}

func (u *UserHomeTablesController) UserHomeTable(DB *gorm.DB) {

	common.NewCreateTableController().Base_Create_Table(DB, &UserHome{}, "user_home")

	if err := seedDefaultHomeInfo(DB); err != nil {
		log.Printf("初始化 user_home 預設資料失敗: %v", err)
	}
}

func seedDefaultHomeInfo(db *gorm.DB) error {
	defaultHomes := []UserHome{
		{Home_id: 1, Community_id: 1, AddressNumber: 1, Floor: "1", Address: "甜水郡社區 1 樓"},
	}

	for _, home := range defaultHomes {
		var existing UserHome
		err := db.Where("home_id = ?", home.Home_id).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if createErr := db.Create(&home).Error; createErr != nil {
				return createErr
			}
			log.Printf("新增預設住戶單位：%d - %s", home.Home_id, home.Address)
		}
	}
	return nil
}
