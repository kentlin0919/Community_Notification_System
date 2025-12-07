package facility_db

import (
	"Community_Notification_System/pkg/databasePkg"

	"gorm.io/gorm"
)

type FacilityController struct{}

func NewFacilityController() *FacilityController {
	return &FacilityController{}
}

func (f *FacilityController) FacilityTable(DB *gorm.DB) {
	databasePkg.NewCreateTableController().Base_Create_Table(DB, &Facilities{}, "facilities")
	databasePkg.NewCreateTableController().Base_Create_Table(DB, &FacilityReservations{}, "facility_reservations")
}
