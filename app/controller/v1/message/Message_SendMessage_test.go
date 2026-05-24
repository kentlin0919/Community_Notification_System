package message

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Community_Notification_System/database"
	message_db "Community_Notification_System/database/Message_DB"
	user_db "Community_Notification_System/database/User_DB"
	"Community_Notification_System/pkg/firebase"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupMessageTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("建立測試資料庫失敗: %v", err)
	}

	if err := db.AutoMigrate(&user_db.UserInfo{}, &message_db.MessageInfo{}); err != nil {
		t.Fatalf("自動遷移資料表失敗: %v", err)
	}

	database.DB = db
	previousFcmClient := firebase.FcmClient
	firebase.FcmClient = nil

	t.Cleanup(func() {
		firebase.FcmClient = previousFcmClient
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		database.DB = nil
	})
}

func TestSendMessageCreatesRecordsWhenFirebaseUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMessageTestDB(t)

	users := []user_db.UserInfo{
		{ID: "user-1", Email: "user1@example.com", Name: "住戶一", Community_id: 100, Fcmtoken: ""},
		{ID: "user-2", Email: "user2@example.com", Name: "住戶二", Community_id: 100, Fcmtoken: ""},
		{ID: "user-3", Email: "other@example.com", Name: "其他社區", Community_id: 200, Fcmtoken: ""},
	}
	if err := database.DB.Create(&users).Error; err != nil {
		t.Fatalf("建立測試使用者失敗: %v", err)
	}

	body := `{
		"title":"社區公告",
		"subtitle":"電梯維修通知",
		"body":"本週三下午 1:00 至 4:00 將進行電梯例行維護。",
		"target_type":"selected_users",
		"recipient_user_ids":["user-1","user-2","user-3"],
		"recipient_emails":["missing@example.com"],
		"metadata":{"category":"announcement","entity_type":"community","entity_id":"notice-1"}
	}`

	controller := NewMessageController()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/messages/send", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Set("user_id", "sender-1")
	ctx.Set("permission_id", 2)
	ctx.Set("community_id", uint64(100))

	controller.SendMessage(ctx)

	if w.Code != http.StatusOK {
		t.Fatalf("預期回傳狀態碼 %d，實際為 %d，body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var response struct {
		Message string `json:"message"`
		Data    struct {
			TargetCount       int      `json:"target_count"`
			SuccessCount      int      `json:"success_count"`
			FailureCount      int      `json:"failure_count"`
			CreatedMessageIDs []string `json:"created_message_ids"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析回應 JSON 失敗: %v", err)
	}
	if response.Data.TargetCount != 2 {
		t.Fatalf("預期只鎖定同社區 2 位收件者，實際為 %d", response.Data.TargetCount)
	}
	if response.Data.SuccessCount != 0 || response.Data.FailureCount != 2 {
		t.Fatalf("Firebase 未初始化時預期 0 成功 2 失敗，實際 success=%d failure=%d", response.Data.SuccessCount, response.Data.FailureCount)
	}
	if len(response.Data.CreatedMessageIDs) != 2 {
		t.Fatalf("預期建立 2 筆 message id，實際為 %d", len(response.Data.CreatedMessageIDs))
	}

	var records []message_db.MessageInfo
	if err := database.DB.Order("user_id asc").Find(&records).Error; err != nil {
		t.Fatalf("查詢 message_info 失敗: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("預期建立 2 筆寄送紀錄，實際為 %d", len(records))
	}
	if records[0].CommunityID != 100 || records[1].CommunityID != 100 {
		t.Fatalf("寄送紀錄 community_id 錯誤: %+v", records)
	}
	if records[0].SenderID != "sender-1" || records[1].SenderID != "sender-1" {
		t.Fatalf("寄送紀錄 sender_id 錯誤: %+v", records)
	}
	if records[0].BatchID == "" || records[0].BatchID != records[1].BatchID {
		t.Fatalf("預期同批次紀錄共用 batch_id，實際為 %+v", records)
	}
	if records[0].FcmStatus != "skipped" || records[1].FcmStatus != "skipped" {
		t.Fatalf("Firebase 未初始化時預期 fcm_status=skipped，實際為 %+v", records)
	}
}
