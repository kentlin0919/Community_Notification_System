package ApiRoute_DB

import (
	"gorm.io/gorm"
)

type ApiRouteController struct{}

func NewApiRouteController() *ApiRouteController {
	return &ApiRouteController{}
}

func (c *ApiRouteController) ApiRouteTable(DB *gorm.DB) {
	DB.AutoMigrate(&ApiRoute{})
}
