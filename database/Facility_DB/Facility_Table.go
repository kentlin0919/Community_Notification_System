package facility_db

import (
	"Community_Notification_System/pkg/common"

	"gorm.io/gorm"
)

type FacilityDBController struct{}

func NewFacilityDBController() *FacilityDBController {
	return &FacilityDBController{}
}

func (u *FacilityDBController) FacilityTable(DB *gorm.DB) {
	common.NewCreateTableController().Base_Create_Table(DB, &FacilityInfo{}, "facility_info")
	common.NewCreateTableController().Base_Create_Table(DB, &FacilityRule{}, "facility_rule")
	common.NewCreateTableController().Base_Create_Table(DB, &FacilityReservation{}, "facility_reservation")
	common.NewCreateTableController().Base_Create_Table(DB, &RescheduleRequest{}, "reschedule_request")
}
