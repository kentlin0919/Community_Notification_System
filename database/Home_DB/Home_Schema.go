package Home_db

type UserHome struct {
	Home_id uint64 `json:"home_id" gorm:"primarykey"`
	Address string `json:"address"`
}
