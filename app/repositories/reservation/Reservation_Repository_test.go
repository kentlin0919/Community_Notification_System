package reservation

import (
	"fmt"
	"testing"

	"Community_Notification_System/database"
	facilitydb "Community_Notification_System/database/Facility_DB"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupReservationRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("建立測試資料庫失敗: %v", err)
	}
	if err := db.AutoMigrate(&facilitydb.FacilityInfo{}, &facilitydb.FacilityRule{}, &facilitydb.FacilityReservation{}); err != nil {
		t.Fatalf("自動遷移資料表失敗: %v", err)
	}

	database.DB = db
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		database.DB = nil
	})
	return db
}

func TestCreateReservationRepositoryRejectsOverlappingActiveReservations(t *testing.T) {
	setupReservationRepositoryTestDB(t)

	existing := facilitydb.FacilityReservation{
		FacilityID:      1,
		UserID:          "resident-1",
		CommunityID:     100,
		ReservationDate: "2026-05-20",
		StartTime:       "10:00",
		EndTime:         "11:00",
		PeopleCount:     2,
		Status:          "approved",
	}
	if err := database.DB.Create(&existing).Error; err != nil {
		t.Fatalf("建立既有預約失敗: %v", err)
	}

	newReservation := &facilitydb.FacilityReservation{
		FacilityID:      1,
		UserID:          "resident-2",
		CommunityID:     100,
		ReservationDate: "2026-05-20",
		StartTime:       "10:30",
		EndTime:         "11:30",
		PeopleCount:     1,
	}
	rule := &facilitydb.FacilityRule{AutoApprove: true}

	result := CreateReservationRepository(newReservation, rule)
	if result.Statue.Error == nil {
		t.Fatalf("重疊時段應被拒絕")
	}

	var count int64
	if err := database.DB.Model(&facilitydb.FacilityReservation{}).Count(&count).Error; err != nil {
		t.Fatalf("查詢預約數量失敗: %v", err)
	}
	if count != 1 {
		t.Fatalf("重疊預約不應被寫入，count=%d", count)
	}
}

func TestGetReservationListRepositoryFiltersResidentsByUser(t *testing.T) {
	setupReservationRepositoryTestDB(t)

	reservations := []facilitydb.FacilityReservation{
		{FacilityID: 1, UserID: "resident-1", CommunityID: 100, ReservationDate: "2026-05-20", StartTime: "09:00", EndTime: "10:00", PeopleCount: 1, Status: "approved"},
		{FacilityID: 1, UserID: "resident-2", CommunityID: 100, ReservationDate: "2026-05-20", StartTime: "10:00", EndTime: "11:00", PeopleCount: 1, Status: "approved"},
		{FacilityID: 2, UserID: "resident-3", CommunityID: 200, ReservationDate: "2026-05-20", StartTime: "11:00", EndTime: "12:00", PeopleCount: 1, Status: "approved"},
	}
	if err := database.DB.Create(&reservations).Error; err != nil {
		t.Fatalf("建立測試預約失敗: %v", err)
	}

	residentResult := GetReservationListRepository(100, "resident-1", true)
	if residentResult.Statue.Error != nil {
		t.Fatalf("住戶查詢預約列表失敗: %v", residentResult.Statue.Error)
	}
	if len(residentResult.Result) != 1 || residentResult.Result[0].UserID != "resident-1" {
		t.Fatalf("住戶應只看到自己的預約，實際=%+v", residentResult.Result)
	}

	adminResult := GetReservationListRepository(100, "admin-1", false)
	if adminResult.Statue.Error != nil {
		t.Fatalf("管理員查詢預約列表失敗: %v", adminResult.Statue.Error)
	}
	if len(adminResult.Result) != 2 {
		t.Fatalf("管理員應看到同社區所有預約，實際=%+v", adminResult.Result)
	}
}

func TestGetReservationDetailRepositoryHonorsRoleScope(t *testing.T) {
	setupReservationRepositoryTestDB(t)

	reservation := facilitydb.FacilityReservation{FacilityID: 1, UserID: "resident-1", CommunityID: 100, ReservationDate: "2026-05-20", StartTime: "09:00", EndTime: "10:00", PeopleCount: 1, Status: "approved"}
	if err := database.DB.Create(&reservation).Error; err != nil {
		t.Fatalf("建立測試預約失敗: %v", err)
	}

	residentResult := GetReservationDetailRepository(reservation.ID, 100, "resident-1", true)
	if residentResult.Statue.Error != nil || residentResult.Result.ID != reservation.ID {
		t.Fatalf("住戶應可查自己的預約詳情，result=%+v err=%v", residentResult.Result, residentResult.Statue.Error)
	}

	otherResidentResult := GetReservationDetailRepository(reservation.ID, 100, "resident-2", true)
	if otherResidentResult.Statue.Error == nil {
		t.Fatalf("住戶不應可查其他人的預約詳情")
	}

	adminResult := GetReservationDetailRepository(reservation.ID, 100, "admin-1", false)
	if adminResult.Statue.Error != nil || adminResult.Result.ID != reservation.ID {
		t.Fatalf("管理員應可查同社區預約詳情，result=%+v err=%v", adminResult.Result, adminResult.Statue.Error)
	}
}

func TestAdminDeleteReservationRepositorySoftDeletesReservation(t *testing.T) {
	setupReservationRepositoryTestDB(t)

	reservation := facilitydb.FacilityReservation{FacilityID: 1, UserID: "resident-1", CommunityID: 100, ReservationDate: "2026-05-20", StartTime: "09:00", EndTime: "10:00", PeopleCount: 1, Status: "approved"}
	if err := database.DB.Create(&reservation).Error; err != nil {
		t.Fatalf("建立測試預約失敗: %v", err)
	}

	if err := DeleteReservationRepository(reservation.ID, 100); err != nil {
		t.Fatalf("管理員軟刪除預約失敗: %v", err)
	}

	var count int64
	if err := database.DB.Model(&facilitydb.FacilityReservation{}).Where("id = ?", reservation.ID).Count(&count).Error; err != nil {
		t.Fatalf("查詢未刪除預約失敗: %v", err)
	}
	if count != 0 {
		t.Fatalf("GORM 預設查詢不應看到已軟刪除預約，count=%d", count)
	}

	var deleted facilitydb.FacilityReservation
	if err := database.DB.Unscoped().First(&deleted, reservation.ID).Error; err != nil {
		t.Fatalf("Unscoped 應查得到已軟刪除預約: %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatalf("DeletedAt 應被設定: %+v", deleted)
	}
}
