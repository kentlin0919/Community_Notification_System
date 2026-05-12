package message_model

type MessageData struct {
	Userselect []string `json:"Userselect" example:"user1@example.com,user2@example.com"`
	IsAllUser  bool     `json:"IsAllUser" example:"false"`
	Title      string   `json:"Title" example:"社區公告"`
	Subtile    string   `json:"Subtile" example:"電梯維修通知"`
	Detail     string   `json:"Detail" example:"本週三下午 1:00 至 4:00 將進行電梯例行維護，請住戶改用樓梯。"`
}

type MessageRequest struct {
	Message string `json:"Message" example:"訊息發送成功"`
}

type FCMNotificationRequest struct {
	DeviceToken string   `json:"deviceToken" binding:"required" example:"fcm-device-token"`
	Title       string   `json:"title" binding:"required" example:"您有新的包裹"`
	Body        string   `json:"body" binding:"required" example:"您的包裹已抵達管理室，請抽空領取。"`
	Userselect  []string `json:"Userselect" example:"user-uuid-1,user-uuid-2"`
}

type MessageListQuery struct {
	Page     int   `form:"page"`
	PageSize int   `form:"page_size"`
	IsRead   *bool `form:"is_read"`
}

type MessageSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Subtile    string `json:"subtile"`
	Detail     string `json:"detail"`
	IsRead     bool   `json:"is_read"`
	CreateTime string `json:"create_time"`
}

type MessageListResponse struct {
	Total    int64            `json:"total"`
	Messages []MessageSummary `json:"messages"`
}
