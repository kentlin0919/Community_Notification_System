package Home_db

type UserHome struct {
	Home_id       uint64 `json:"home_id" gorm:"primarykey"`
	Community_id  uint64 `json:"community_id"`
	AddressNumber int    `json:"AddressNumber"`
	Floor         string `json:"floor"`
	Address       string `json:"address"`
}
