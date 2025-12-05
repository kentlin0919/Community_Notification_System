package database

import (
	facility_db "Community_Notification_System/database/Facility_DB"
	homedb "Community_Notification_System/database/Home_DB"
	maintenance_db "Community_Notification_System/database/Maintenance_DB"
	platform_db "Community_Notification_System/database/Platform_DB"
	visitor_db "Community_Notification_System/database/Visitor_DB"

	message_db "Community_Notification_System/database/Message_DB"
	permission_db "Community_Notification_System/database/Permission_DB"
	userlog_db "Community_Notification_System/database/UserLog_DB"
	user_db "Community_Notification_System/database/User_DB"

	communitydb "Community_Notification_System/database/Community_DB"

	actionlog_db "Community_Notification_System/database/ActionLog_DB"
	apiroute_db "Community_Notification_System/database/ApiRoute_DB"

	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// 取得環境變數
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	// 建立 DSN 連線字串
	dsn := buildDSN(dbHost, dbUser, dbPassword, dbName, dbPort, dbTimezone)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		if isMissingDatabaseError(err) {
			log.Printf("偵測到資料庫 %s 不存在，嘗試自動建立...", dbName)
			if createErr := createDatabaseIfNotExists(dbHost, dbUser, dbPassword, dbName, dbPort, dbTimezone); createErr != nil {
				log.Fatalf("建立資料庫 %s 失敗: %v", dbName, createErr)
			}
			DB, err = gorm.Open(postgres.New(postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true,
			}), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})
		}
		if err != nil {
			log.Fatalf("資料庫連線失敗: %v", err)
		}
	}
	log.Println("資料庫連線成功")
	CreateTable()
}

func buildDSN(host, user, password, dbName, port, timezone string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, user, password, dbName, port, timezone,
	)
}

func isMissingDatabaseError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "does not exist") && strings.Contains(errMsg, "database")
}

func createDatabaseIfNotExists(host, user, password, dbName, port, timezone string) error {
	adminDSN := buildDSN(host, user, password, "postgres", port, timezone)
	adminDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  adminDSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("連線至預設資料庫失敗: %w", err)
	}

	sqlDB, err := adminDB.DB()
	if err != nil {
		return fmt.Errorf("取得資料庫連線失敗: %w", err)
	}
	defer sqlDB.Close()

	createSQL := fmt.Sprintf("CREATE DATABASE %q", dbName)
	if execErr := adminDB.Exec(createSQL).Error; execErr != nil {
		if strings.Contains(strings.ToLower(execErr.Error()), "already exists") {
			return nil
		}
		return fmt.Errorf("執行建庫指令失敗: %w", execErr)
	}
	log.Printf("資料庫 %s 已建立。", dbName)
	return nil
}

func CreateTable() {
	/// 創建UserInfo
	user_db.NewUserDBController().UserTable(DB)

	/// 創建User Home
	homedb.NewUserHomeTableController().UserHomeTable(DB)

	/// 創建User Log
	userlog_db.NewUserLogTableController().UserLogTable(DB)

	/// 創建Message Table
	message_db.NewUserDBController().MessageInfoTable(DB)

	/// 創建Permission table
	permission_db.NewPermissionInfoController().PermissionInfoTable(DB)

	/// 創建Platform_DB
	platform_db.NewPlatformInfoController().PlatformInfoTable(DB)

	/// 創建Community_DB
	communitydb.NewCommunityInfoController().CommunityInfoTable(DB)

	/// 創建ActionLog_DB
	actionlog_db.NewActionLogController().ActionLogTable(DB)

	/// 創建ApiRoute_DB
	apiroute_db.NewApiRouteController().ApiRouteTable(DB)

	/// 創建Facility_DB
	facility_db.NewFacilityController().FacilityTable(DB)

	/// 創建Visitor_DB
	visitor_db.NewVisitorController().VisitorTable(DB)

	/// 創建Maintenance_DB
	maintenance_db.NewMaintenanceController().MaintenanceTable(DB)
}
