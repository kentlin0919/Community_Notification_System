package facility

import (
	"fmt"
	"testing"

	"Community_Notification_System/database"
	facilitydb "Community_Notification_System/database/Facility_DB"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupFacilityRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("建立測試資料庫失敗: %v", err)
	}
	if err := db.AutoMigrate(&facilitydb.FacilityInfo{}, &facilitydb.FacilityRule{}); err != nil {
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

func TestGetFacilityListRepositoryFiltersActiveForResidents(t *testing.T) {
	setupFacilityRepositoryTestDB(t)

	facilities := []facilitydb.FacilityInfo{
		{CommunityID: 100, Name: "健身房", FacilityType: "sport", Status: "active"},
		{CommunityID: 100, Name: "會議室", FacilityType: "meeting", Status: "inactive"},
		{CommunityID: 200, Name: "其他社區健身房", FacilityType: "sport", Status: "active"},
	}
	if err := database.DB.Create(&facilities).Error; err != nil {
		t.Fatalf("建立測試設施失敗: %v", err)
	}

	residentResult := GetFacilityListRepository(100, true)
	if residentResult.Statue.Error != nil {
		t.Fatalf("住戶查詢設施列表失敗: %v", residentResult.Statue.Error)
	}
	if len(residentResult.Result) != 1 || residentResult.Result[0].Name != "健身房" {
		t.Fatalf("住戶應只看到同社區 active 設施，實際=%+v", residentResult.Result)
	}

	adminResult := GetFacilityListRepository(100, false)
	if adminResult.Statue.Error != nil {
		t.Fatalf("管理員查詢設施列表失敗: %v", adminResult.Statue.Error)
	}
	if len(adminResult.Result) != 2 {
		t.Fatalf("管理員應看到同社區所有未刪除設施，實際=%+v", adminResult.Result)
	}
}

func TestUpdateFacilityRepositoryUpdatesFacilityAndRule(t *testing.T) {
	setupFacilityRepositoryTestDB(t)

	facility := facilitydb.FacilityInfo{CommunityID: 100, Name: "交誼廳", FacilityType: "room", Status: "draft"}
	if err := database.DB.Create(&facility).Error; err != nil {
		t.Fatalf("建立測試設施失敗: %v", err)
	}
	rule := facilitydb.FacilityRule{FacilityID: facility.ID, SlotMinutes: 60, MaxAdvanceDays: 7, CancelBeforeHours: 24, AutoApprove: false, IsActive: false}
	if err := database.DB.Create(&rule).Error; err != nil {
		t.Fatalf("建立測試設施規則失敗: %v", err)
	}

	facility.Name = "多功能教室"
	facility.Status = "active"
	rule.SlotMinutes = 30
	rule.MaxAdvanceDays = 14
	rule.CancelBeforeHours = 12
	rule.AutoApprove = true
	rule.IsActive = true

	result := UpdateFacilityWithRuleRepository(&facility, &rule)
	if result.Statue.Error != nil {
		t.Fatalf("更新設施與規則失敗: %v", result.Statue.Error)
	}

	var updatedFacility facilitydb.FacilityInfo
	if err := database.DB.First(&updatedFacility, facility.ID).Error; err != nil {
		t.Fatalf("查詢更新後設施失敗: %v", err)
	}
	if updatedFacility.Name != "多功能教室" || updatedFacility.Status != "active" {
		t.Fatalf("設施未正確更新: %+v", updatedFacility)
	}

	var updatedRule facilitydb.FacilityRule
	if err := database.DB.Where("facility_id = ?", facility.ID).First(&updatedRule).Error; err != nil {
		t.Fatalf("查詢更新後規則失敗: %v", err)
	}
	if updatedRule.SlotMinutes != 30 || updatedRule.MaxAdvanceDays != 14 || updatedRule.CancelBeforeHours != 12 || !updatedRule.AutoApprove || !updatedRule.IsActive {
		t.Fatalf("設施規則未正確同步更新: %+v", updatedRule)
	}
}

func TestDeleteFacilityRepositorySoftDeletesFacility(t *testing.T) {
	setupFacilityRepositoryTestDB(t)

	facility := facilitydb.FacilityInfo{CommunityID: 100, Name: "游泳池", FacilityType: "pool", Status: "active"}
	if err := database.DB.Create(&facility).Error; err != nil {
		t.Fatalf("建立測試設施失敗: %v", err)
	}

	if err := DeleteFacilityRepository(facility.ID, 100); err != nil {
		t.Fatalf("軟刪除設施失敗: %v", err)
	}

	var count int64
	if err := database.DB.Model(&facilitydb.FacilityInfo{}).Where("id = ?", facility.ID).Count(&count).Error; err != nil {
		t.Fatalf("查詢未刪除設施失敗: %v", err)
	}
	if count != 0 {
		t.Fatalf("GORM 預設查詢不應看到已軟刪除設施，count=%d", count)
	}

	var deleted facilitydb.FacilityInfo
	if err := database.DB.Unscoped().First(&deleted, facility.ID).Error; err != nil {
		t.Fatalf("Unscoped 應查得到已軟刪除設施: %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatalf("DeletedAt 應被設定: %+v", deleted)
	}
}
