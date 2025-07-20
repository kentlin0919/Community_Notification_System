package model

type ErrorRequest struct {
	Error string `json:"error" example:"error"`
}

type RequestMessage struct {
	Message string `json:"Message" example:"message"`
}
