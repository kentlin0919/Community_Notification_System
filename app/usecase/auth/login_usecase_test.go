package auth

import (
	"os"
	"testing"

	accountModel "Community_Notification_System/app/models/account"
	"Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"
	userlog_db "Community_Notification_System/database/UserLog_DB"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupLoginUsecaseTestDB(t *testing.T) {
	t.Helper()
	t.Setenv("JWTPASSWORD", "testsecret1234567890")

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("建立測試資料庫失敗: %v", err)
	}
	if err := db.AutoMigrate(&user_db.UserInfo{}, &user_db.UserSession{}, &userlog_db.UserLog{}); err != nil {
		t.Fatalf("自動遷移資料表失敗: %v", err)
	}
	database.DB = db
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		database.DB = nil
	})
	_ = os.Getenv("JWTPASSWORD")
}

func TestExecuteLoginCreatesSessionAndReturnsRefreshToken(t *testing.T) {
	setupLoginUsecaseTestDB(t)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	database.DB.Create(&user_db.UserInfo{ID: "user-1", Email: "user@example.com", Password: string(hashedPassword)})

	result, err := ExecuteLogin(&accountModel.User{Email: "user@example.com", Password: "password123", Platform: "App"})
	if err != nil {
		t.Fatalf("預期登入成功，實際錯誤: %v", err)
	}
	if result.RefreshToken == "" {
		t.Fatal("預期回傳非空的 RefreshToken")
	}

	var session user_db.UserSession
	if err := database.DB.Where("user_id = ? AND refresh_token = ?", "user-1", result.RefreshToken).First(&session).Error; err != nil {
		t.Fatalf("預期建立對應的 UserSession，實際查詢失敗: %v", err)
	}
	if session.IsRevoked {
		t.Fatal("預期新建立的 session 未被撤銷")
	}
}
