package user_db

type UserInfo struct {
	ID       uint64 `json:"id" gorm:"primarykey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Ssotoken string `json:"Ssotoken"`
	Home_id  string `json:"Home_id"`
}
