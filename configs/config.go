package configs

import (
	"log"

	"github.com/joho/godotenv"
)

func InitConfig() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("無法載入 .env 文件: %v", err)
	}
}
