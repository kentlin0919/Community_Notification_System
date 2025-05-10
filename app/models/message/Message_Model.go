package message_model

type MessageData struct {
	Title   string `json:"Title" example:"Test"`
	Subtile string `json:"Subtile" example:"Subtile"`
	Detail  string `json:"Detail" example:"Detail"`
}

type MessageRequest struct {
	Message string `json:"Message" example:"Sucessful send Message"`
}
