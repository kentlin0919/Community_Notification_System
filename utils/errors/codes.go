package errors

const (
	// Success
	Success = 0

	// Common Errors (1000-1999)
	ErrInternal      = 1000
	ErrInvalidParams = 1001
	ErrUnauthorized  = 1002
	ErrForbidden     = 1003
	ErrNotFound      = 1004
	ErrConflict      = 1005
	ErrDatabase      = 1006

	// User Related (2000-2999)
	ErrUserNotFound      = 2001
	ErrPasswordWrong     = 2002
	ErrTokenInvalid      = 2003
	ErrOtpInvalid        = 2004
	ErrOtpLocked         = 2005
	ErrResetTokenInvalid = 2006

	// Community Related (3000-3999)
	ErrCommunityNotFound = 3001
	ErrNoPermission      = 3002

	// IoT/Smart Access Related (4000-4999)
	ErrIoTDeviceOffline = 4001
	ErrInvalidPassCode  = 4002

	// Maintenance Related (5000-5999)
	ErrMaintenanceNotFound = 5001
	ErrAssignedFailed      = 5002
)

var ErrorMsg = map[int]string{
	Success:                "操作成功",
	ErrInternal:            "伺服器內部錯誤",
	ErrInvalidParams:       "參數錯誤",
	ErrUnauthorized:        "未經授權",
	ErrForbidden:           "拒絕執行",
	ErrNotFound:            "找不到目標",
	ErrConflict:            "資料衝突",
	ErrDatabase:            "資料庫操作失敗",
	ErrUserNotFound:        "查無此使用者",
	ErrPasswordWrong:       "密碼錯誤",
	ErrTokenInvalid:        "無效的授權 token",
	ErrOtpInvalid:          "驗證碼錯誤或已過期",
	ErrOtpLocked:           "嘗試次數過多，請重新申請驗證碼",
	ErrResetTokenInvalid:   "無效或已過期的重設密碼權杖",
	ErrCommunityNotFound:   "查無此社區",
	ErrNoPermission:        "權限不足",
	ErrIoTDeviceOffline:    "物聯網設備離線",
	ErrInvalidPassCode:     "無效的通行碼",
	ErrMaintenanceNotFound: "查無此維修單",
	ErrAssignedFailed:      "派工失敗",
}

// GetMsg 取得代碼對應的訊息
func GetMsg(code int) string {
	if msg, ok := ErrorMsg[code]; ok {
		return msg
	}
	return ErrorMsg[ErrInternal]
}
