package facility_db

import (
	"Community_Notification_System/pkg/databasepkg"

	"gorm.io/gorm"
)

type FacilityController struct{}

func NewFacilityController() *FacilityController {
	return &FacilityController{}
}

func (f *FacilityController) FacilityTable(DB *gorm.DB) {
	databasepkg.NewCreateTableController().Base_Create_Table(DB, &Facilities{}, "facilities")
	databasepkg.NewCreateTableController().Base_Create_Table(DB, &FacilityReservations{}, "facility_reservations")
}
