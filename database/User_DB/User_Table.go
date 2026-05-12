package user_db

import (
	"Community_Notification_System/pkg/common"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserTablesController struct{}

func NewUserDBController() *UserTablesController {
	return &UserTablesController{}
}

func (u *UserTablesController) UserTable(DB *gorm.DB) {
	// 檢查是否存在 UserInfo 表

	common.NewCreateTableController().Base_Create_Table(DB, &UserInfo{}, "user_info")
	if err := seedDefaultUsers(DB); err != nil {
		log.Printf("初始化 UserInfo 預設資料失敗: %v", err)
	}
}

func seedDefaultUsers(db *gorm.DB) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("09190919"), bcrypt.DefaultCost)

	defaultUsers := []UserInfo{
		{
			ID:           uuid.New().String(),
			Email:        "kent900919@gmail.com",
			Password:     string(hashedPassword),
			Name:         "系統管理員",
			PermissionId: 1,
			Platform:     3,
			Community_id: 0,
			Birthdaytime: time.Now(),
			Registertime: time.Now(),
		},
		{
			ID:           uuid.New().String(),
			Email:        "super@example.com",
			Password:     string(hashedPassword),
			Name:         "超級管理員",
			PermissionId: 1,
			Platform:     3,
			Community_id: 0,
			Birthdaytime: time.Now(),
			Registertime: time.Now(),
		},
		{
			ID:           uuid.New().String(),
			Email:        "admin@example.com",
			Password:     string(hashedPassword),
			Name:         "社區管理員",
			PermissionId: 2,
			Platform:     3,
			Community_id: 1,
			Birthdaytime: time.Now(),
			Registertime: time.Now(),
		},
		{
			ID:           uuid.New().String(),
			Email:        "staff@example.com",
			Password:     string(hashedPassword),
			Name:         "保全人員",
			PermissionId: 3,
			Platform:     3,
			Community_id: 1,
			Birthdaytime: time.Now(),
			Registertime: time.Now(),
		},
		{
			ID:           uuid.New().String(),
			Email:        "resident@example.com",
			Password:     string(hashedPassword),
			Name:         "測試住戶",
			PermissionId: 8,
			Platform:     2, // App
			Community_id: 1,
			Home_id:      "1",
			Birthdaytime: time.Now(),
			Registertime: time.Now(),
		},
	}

	for _, user := range defaultUsers {
		var existing UserInfo
		err := db.Where("Email = ?", user.Email).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if createErr := db.Create(&user).Error; createErr != nil {
				return fmt.Errorf("新增預設使用者 %s 失敗: %w", user.Email, createErr)
			}
			log.Printf("新增預設使用者：%s - %s (權限: %d)", user.Email, user.Name, user.PermissionId)
		}
	}

	return nil
}
