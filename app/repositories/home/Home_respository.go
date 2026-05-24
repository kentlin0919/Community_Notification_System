package repository

import (
	"Community_Notification_System/database"
	communitydb "Community_Notification_System/database/Community_DB"
	facilitydb "Community_Notification_System/database/Facility_DB"
	messagedb "Community_Notification_System/database/Message_DB"
	parceldb "Community_Notification_System/database/Parcel_DB"
	userdb "Community_Notification_System/database/User_DB"
	"time"
)

// ── 住戶查詢 ──

// CountUnreadMessages 計算使用者未讀通知數
func CountUnreadMessages(userID string) int64 {
	var count int64
	database.DB.Model(&messagedb.MessageInfo{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count)
	return count
}

// CountUnclaimedParcels 計算社區待領包裹數（status=1 代表待領取）
func CountUnclaimedParcels(communityID uint64) int64 {
	var count int64
	database.DB.Model(&parceldb.ParcelInfo{}).
		Where("community_id = ? AND status = ?", communityID, 1).
		Count(&count)
	return count
}

// CountUpcomingReservations 計算使用者未來預約數
func CountUpcomingReservations(userID string) int64 {
	var count int64
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&facilitydb.FacilityReservation{}).
		Where("user_id = ? AND reservation_date >= ? AND status IN ?", userID, today, []string{"pending", "approved"}).
		Count(&count)
	return count
}

// GetRecentMessages 取得使用者最新 N 筆通知
func GetRecentMessages(userID string, limit int) []messagedb.MessageInfo {
	var messages []messagedb.MessageInfo
	database.DB.Where("user_id = ?", userID).
		Order("create_time DESC").
		Limit(limit).
		Find(&messages)
	return messages
}

// ReservationWithFacility 帶設施名稱的預約（首頁聚合用）
type ReservationWithFacility struct {
	ID              uint64 `json:"id"`
	FacilityName    string `json:"facility_name"`
	ReservationDate string `json:"reservation_date"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	Status          string `json:"status"`
}

// GetUpcomingReservations 取得使用者最近 N 筆未來預約（含設施名稱）
func GetUpcomingReservations(userID string, limit int) []ReservationWithFacility {
	var results []ReservationWithFacility
	today := time.Now().Format("2006-01-02")
	database.DB.Table("facility_reservation").
		Select("facility_reservation.id, facility_info.name as facility_name, facility_reservation.reservation_date, facility_reservation.start_time, facility_reservation.end_time, facility_reservation.status").
		Joins("LEFT JOIN facility_info ON facility_info.id = facility_reservation.facility_id").
		Where("facility_reservation.user_id = ? AND facility_reservation.reservation_date >= ? AND facility_reservation.status IN ? AND facility_reservation.deleted_at IS NULL",
			userID, today, []string{"pending", "approved"}).
		Order("facility_reservation.reservation_date ASC, facility_reservation.start_time ASC").
		Limit(limit).
		Scan(&results)
	return results
}

// ── 社區管理員查詢 ──

// CountCommunityResidents 計算社區住戶數
func CountCommunityResidents(communityID uint64) int64 {
	var count int64
	database.DB.Model(&userdb.UserInfo{}).
		Where("community_id = ?", communityID).
		Count(&count)
	return count
}

// CountTodayReservations 計算社區今日預約數
func CountTodayReservations(communityID uint64) int64 {
	var count int64
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&facilitydb.FacilityReservation{}).
		Where("community_id = ? AND reservation_date = ?", communityID, today).
		Count(&count)
	return count
}

// GetTodayReservations 取得社區今日的預約列表
func GetTodayReservations(communityID uint64, limit int) []ReservationWithFacility {
	var results []ReservationWithFacility
	today := time.Now().Format("2006-01-02")
	database.DB.Table("facility_reservation").
		Select("facility_reservation.id, facility_info.name as facility_name, facility_reservation.reservation_date, facility_reservation.start_time, facility_reservation.end_time, facility_reservation.status").
		Joins("LEFT JOIN facility_info ON facility_info.id = facility_reservation.facility_id").
		Where("facility_reservation.community_id = ? AND facility_reservation.reservation_date = ? AND facility_reservation.deleted_at IS NULL",
			communityID, today).
		Order("facility_reservation.start_time ASC").
		Limit(limit).
		Scan(&results)
	return results
}

// GetCommunityRecentMessages 取得社區最新 N 筆通知
func GetCommunityRecentMessages(communityID uint64, limit int) []messagedb.MessageInfo {
	var messages []messagedb.MessageInfo
	database.DB.Where("community_id = ?", communityID).
		Order("create_time DESC").
		Limit(limit).
		Find(&messages)
	return messages
}

// ── 超級管理員查詢 ──

// CountAllCommunities 計算全平台社區總數
func CountAllCommunities() int64 {
	var count int64
	database.DB.Model(&communitydb.CommunityInfo{}).Count(&count)
	return count
}

// CountAllUsers 計算全平台使用者總數
func CountAllUsers() int64 {
	var count int64
	database.DB.Model(&userdb.UserInfo{}).Count(&count)
	return count
}

// CountPendingCommunityApplications 計算待審核社區申請數
func CountPendingCommunityApplications() int64 {
	var count int64
	database.DB.Model(&communitydb.CommunityRegisterApplication{}).
		Where("status = ?", "pending").
		Count(&count)
	return count
}

// GetPendingCommunityApplications 取得待審核社區申請列表
func GetPendingCommunityApplications(limit int) []communitydb.CommunityRegisterApplication {
	var apps []communitydb.CommunityRegisterApplication
	database.DB.Where("status = ?", "pending").
		Order("created_at DESC").
		Limit(limit).
		Find(&apps)
	return apps
}
