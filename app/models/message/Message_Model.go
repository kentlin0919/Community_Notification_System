package message_model

type MessageData struct {
	Userselect []string `json:"Userselect" example:"user1@example.com,user2@example.com"`
	IsAllUser  bool     `json:"IsAllUser" example:"false"`
	Title      string   `json:"Title" example:"Test"`
	Subtile    string   `json:"Subtile" example:"Subtile"`
	Detail     string   `json:"Detail" example:"Detail"`
}

type MessageRequest struct {
	Message string `json:"Message" example:"Sucessful send Message"`
}
