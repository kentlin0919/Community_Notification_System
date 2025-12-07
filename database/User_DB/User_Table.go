package user_db

import (
	"Community_Notification_System/app/models/model"
	"Community_Notification_System/pkg/databasePkg"
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

	databasePkg.NewCreateTableController().Base_Create_Table(DB, &UserInfo{}, "user_info")
	if err := seedDefaultAdminUser(DB); err != nil {
		log.Printf("初始化 UserInfo 預設資料失敗: %v", err)
	}
}

func seedDefaultAdminUser(db *gorm.DB) error {

	hashedPassword, errHashedPassword := bcrypt.GenerateFromPassword([]byte("09190919"), bcrypt.DefaultCost)
	if errHashedPassword != nil {
		var errorModel model.ErrorRequest
		errorModel.Error = "密碼加密失敗"

	}

	uuidString := uuid.New()
	defaultUser := UserInfo{
		ID:           uuidString.String(),
		Email:        "kent900919@gmail.com",
		Password:     string(hashedPassword), // 存儲加密後的密碼
		BirthdayTime: time.Now(),
		RegisterTime: time.Now(),
		PermissionId: 1,
		Name:         "系統管理員",
		Platform:     3,
		Community_id: 0,
	}

	var existing UserInfo
	err := db.Where("Email = ?", defaultUser.Email).First(&existing).Error

	switch {
	case err == nil:

	case errors.Is(err, gorm.ErrRecordNotFound):

		if createErr := db.Create(&defaultUser).Error; createErr != nil {
			return fmt.Errorf("新增預設Admin %s 失敗: %w", defaultUser.Email, createErr)
		}
		log.Printf("新增預設系統管理員：%s - %s", defaultUser.Email, defaultUser.Name)
	default:
		return fmt.Errorf("查詢權限 %s 失敗: %w", defaultUser.Name, err)
	}

	return nil
}
