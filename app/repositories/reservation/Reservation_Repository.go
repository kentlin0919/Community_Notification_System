package reservation

import (
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	facilitydb "Community_Notification_System/database/Facility_DB"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CheckConflict 確認是否有預約時段衝突 (只看 pending 或 approved)
func CheckConflict(tx *gorm.DB, facilityID uint64, date string, start string, end string) error {
	var count int64
	// 結束時間大於新開始時間，且開始時間小於新結束時間 => 有重疊
	// 前提：status 屬於有效預約
	err := tx.Model(&facilitydb.FacilityReservation{}).
		Where("facility_id = ?", facilityID).
		Where("reservation_date = ?", date).
		Where("status IN ?", []string{"pending", "approved"}).
		Where("start_time < ? AND end_time > ?", end, start).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("該時段已有其他預約，發生衝突")
	}
	return nil
}

// CreateReservationRepository 建立一筆預約單
func CreateReservationRepository(req *facilitydb.FacilityReservation, facilityRule *facilitydb.FacilityRule) repositoryModels.RepositoryModel[facilitydb.FacilityReservation] {
	var result repositoryModels.RepositoryModel[facilitydb.FacilityReservation]

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 檢查時段衝突
		if err := CheckConflict(tx, req.FacilityID, req.ReservationDate, req.StartTime, req.EndTime); err != nil {
			return err
		}

		// 2. 判斷審核條件
		if facilityRule.AutoApprove {
			req.Status = "approved"
		} else {
			req.Status = "pending"
		}

		// 3. 建立
		if err := tx.Create(req).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		result.Statue.Error = err
	} else {
		result.Result = *req
		result.Statue.RowsAffected = 1
	}

	return result
}

// CancelReservationRepository 標記預約為取消
func CancelReservationRepository(resID uint64, userID string, reason string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var res facilitydb.FacilityReservation
		if err := tx.First(&res, resID).Error; err != nil {
			return err
		}

		if res.UserID != userID {
			return errors.New("您沒有權限取消此預約")
		}

		if res.Status == "cancelled" || res.Status == "completed" || res.Status == "no_show" {
			return errors.New("該預約當前狀態不可取消")
		}

		return tx.Model(&res).Updates(map[string]interface{}{
			"status":        "cancelled",
			"cancel_reason": reason,
			"updated_at":    time.Now(),
		}).Error
	})
}

// CreateRescheduleRequestRepository 建立改期申請單
func CreateRescheduleRequestRepository(req *facilitydb.RescheduleRequest, userID string, facilityID uint64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 驗證原預約是否存在且屬於該 user
		var res facilitydb.FacilityReservation
		if err := tx.First(&res, req.ReservationID).Error; err != nil {
			return err
		}
		
		if res.UserID != userID || res.FacilityID != facilityID {
			return errors.New("無效的預約操作")
		}

		// 檢查新時段是否衝突
		if err := CheckConflict(tx, facilityID, req.NewDate, req.NewStartTime, req.NewEndTime); err != nil {
			return err
		}

		// 寫入改期申請
		return tx.Create(req).Error
	})
}

// ApproveRescheduleRequestRepository 同意/拒絕改期
func ApproveRescheduleRequestRepository(reqID uint64, isApprove bool, comment string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var reschedule facilitydb.RescheduleRequest
		if err := tx.First(&reschedule, reqID).Error; err != nil {
			return err
		}

		if reschedule.Status != "pending" {
			return errors.New("此改期申請已處理過")
		}

		newStatus := "rejected"
		if isApprove {
			newStatus = "approved"

			// 同意的話要一併查出原預約，並更新它
			var res facilitydb.FacilityReservation
			if err := tx.First(&res, reschedule.ReservationID).Error; err != nil {
				return err
			}

			// 再次檢查是否有衝突 (因為別人可能在此期間搶走時段)
			if err := CheckConflict(tx, res.FacilityID, reschedule.NewDate, reschedule.NewStartTime, reschedule.NewEndTime); err != nil {
				return err
			}

			if err := tx.Model(&res).Updates(map[string]interface{}{
				"reservation_date": reschedule.NewDate,
				"start_time":       reschedule.NewStartTime,
				"end_time":         reschedule.NewEndTime,
			}).Error; err != nil {
				return err
			}
		}

		// 更新改期單的狀態
		return tx.Model(&reschedule).Updates(map[string]interface{}{
			"status":        newStatus,
			"admin_comment": comment,
			"updated_at":    time.Now(),
		}).Error
	})
}
