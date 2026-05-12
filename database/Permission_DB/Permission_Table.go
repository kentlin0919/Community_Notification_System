package permission_db

import (
	"Community_Notification_System/pkg/common"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type PermissionInfoController struct{}

func NewPermissionInfoController() *PermissionInfoController {
	return &PermissionInfoController{}
}

func (u *PermissionInfoController) PermissionInfoTable(DB *gorm.DB) {
	// 檢查是否存在 UserInfo 表
	common.NewCreateTableController().Base_Create_Table(DB, &PermissionInfo{}, "permission_info")
	
	// 建立社區自訂權限表
	common.NewCreateTableController().Base_Create_Table(DB, &CommunityPermissionProfile{}, "community_permission_profile")
	
	if err := seedDefaultPermissions(DB); err != nil {
		log.Printf("初始化 permission_info 預設資料失敗: %v", err)
	}
}

func seedDefaultPermissions(db *gorm.DB) error {
	defaultPermissions := []PermissionInfo{
		{PermissionID: "1", Name: "系統管理員 (Super Admin)"},
		{PermissionID: "2", Name: "社區管理員 (Admin)"},
		{PermissionID: "3", Name: "自訂角色 A (如: 保全)"},
		{PermissionID: "4", Name: "自訂角色 B (如: 主委)"},
		{PermissionID: "5", Name: "自訂角色 C (如: 委員)"},
		{PermissionID: "6", Name: "自訂角色 D (如: 櫃台)"},
		{PermissionID: "7", Name: "自訂角色 E (如: 財務)"},
		{PermissionID: "8", Name: "住戶 (Resident)"},
	}

	for _, perm := range defaultPermissions {
		var existing PermissionInfo
		err := db.Where("permission_id = ?", perm.PermissionID).First(&existing).Error
		switch {
		case err == nil:
			continue
		case errors.Is(err, gorm.ErrRecordNotFound):
			perm.ID = perm.PermissionID
			if createErr := db.Create(&perm).Error; createErr != nil {
				return fmt.Errorf("新增預設權限 %s 失敗: %w", perm.PermissionID, createErr)
			}
			log.Printf("新增預設權限：%s - %s", perm.PermissionID, perm.Name)
		default:
			return fmt.Errorf("查詢權限 %s 失敗: %w", perm.PermissionID, err)
		}
	}
	return nil
}
