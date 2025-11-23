package Home_db

import (
	"Community_Notification_System/pkg/databasepkg"

	"gorm.io/gorm"
)

type UserHomeTablesController struct{}

func NewUserHomeTableController() *UserHomeTablesController {
	return &UserHomeTablesController{}
}

func (u *UserHomeTablesController) UserHomeTable(DB *gorm.DB) {

	databasepkg.NewCreateTableController().Base_Create_Table(DB, &UserHome{}, "user_home")

}
