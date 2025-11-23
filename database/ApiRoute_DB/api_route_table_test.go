package ApiRoute_DB

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestApiRouteTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in‑memory DB: %v", err)
	}

	controller := NewApiRouteController()
	controller.ApiRouteTable(db)

	if !db.Migrator().HasTable(&ApiRoute{}) {
		t.Fatalf("api_routes table was not created")
	}
}
