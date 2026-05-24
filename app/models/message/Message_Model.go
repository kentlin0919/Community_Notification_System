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

type FCMNotificationRequest = SendMessageRequest

type SendMessageRequest struct {
	DeviceToken      string              `json:"deviceToken" example:"fcm-device-token"`
	Title            string              `json:"title" binding:"required" example:"您有新的包裹"`
	Subtitle         string              `json:"subtitle" example:"包裹到達通知"`
	Body             string              `json:"body" binding:"required" example:"您的包裹已抵達管理室，請抽空領取。"`
	TargetType       string              `json:"target_type" example:"selected_users"`
	RecipientUserIDs []string            `json:"recipient_user_ids" example:"user-uuid-1,user-uuid-2"`
	RecipientEmails  []string            `json:"recipient_emails" example:"user1@example.com,user2@example.com"`
	Userselect       []string            `json:"Userselect" example:"user1@example.com,user2@example.com"`
	Metadata         SendMessageMetadata `json:"metadata"`
}

type SendMessageMetadata struct {
	Category    string `json:"category" example:"announcement"`
	EntityType  string `json:"entity_type" example:"community"`
	EntityID    string `json:"entity_id" example:"notice-1"`
	ClickAction string `json:"click_action" example:"FLUTTER_NOTIFICATION_CLICK"`
}

type SendMessageResponse struct {
	Message string                  `json:"message"`
	Data    SendMessageResponseData `json:"data"`
}

type SendMessageResponseData struct {
	MessageBatchID    string              `json:"message_batch_id"`
	TargetCount       int                 `json:"target_count"`
	SuccessCount      int                 `json:"success_count"`
	FailureCount      int                 `json:"failure_count"`
	CreatedMessageIDs []string            `json:"created_message_ids"`
	FCMResults        []SendMessageResult `json:"fcm_results"`
}

type SendMessageResult struct {
	MessageID    string `json:"message_id"`
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Success      bool   `json:"success"`
	FcmStatus    string `json:"fcm_status"`
	FcmMessageID string `json:"fcm_message_id,omitempty"`
	Error        string `json:"error,omitempty"`
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
