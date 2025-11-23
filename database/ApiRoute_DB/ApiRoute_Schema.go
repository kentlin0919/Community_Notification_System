package ApiRoute_DB

import (
	"gorm.io/gorm"
)

type ApiRoute struct {
	gorm.Model
	Path                 string `gorm:"type:varchar(255);uniqueIndex:idx_path_method;comment:API路徑"`
	Method               string `gorm:"type:varchar(10);uniqueIndex:idx_path_method;comment:請求方法"`
	Description          string `gorm:"type:varchar(255);comment:功能描述"`
	RequiredPermissionID uint   `gorm:"comment:所需權限ID"`
}

func (ApiRoute) TableName() string {
	return "api_routes"
}
