package database

import (
	homedb "Community_Notification_System/database/Home_DB"
	userlog_db "Community_Notification_System/database/UserLog_DB"
	user_db "Community_Notification_System/database/User_DB"
	"fmt"
	"log"
	"os"

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
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbTimezone,
	)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatalf("資料庫連線失敗: %v", err)
	}
	log.Println("資料庫連線成功")
	CreateTable()
}

func CreateTable() {
	/// 創建UserInfo
	user_db.NewUserDBController().UserTable(DB)

	/// 創建User Home
	homedb.NewUserHomeTableController().UserHomeTable(DB)

	/// 創建User Log
	userlog_db.NewUserLogTableController().UserLogTable(DB)

}
