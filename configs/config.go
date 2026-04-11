package configs

import (
	"log"

	"github.com/joho/godotenv"
)

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("提示：未發現 .env 文件或預載入失敗，將使用系統環境變數 (%v)", err)
	}
}
