package databasePkg

import (
	"fmt"

	"gorm.io/gorm"
)

type CreateTableController struct{}

func NewCreateTableController() *CreateTableController {
	return &CreateTableController{}
}

func (c *CreateTableController) Base_Create_Table(DB *gorm.DB, schema interface{}, tableName string) {
	// 檢查是否存在 UserInfo 表
	if DB.Migrator().HasTable(tableName) {
		fmt.Printf("%s 表存在\n", tableName)
	} else {
		fmt.Printf("%s 表不存在\n", tableName)

		// 創建 UserInfo 表
		if err := DB.AutoMigrate(schema); err != nil {
			fmt.Printf("創建 %s 表失敗: %v\n", tableName, err)
		} else {
			fmt.Printf("%s 表創建成功\n", tableName)
		}
	}
}
