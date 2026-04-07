package request

type CreateChatRequest struct {
	Message *string `json:"message" binding:"omitempty"`
}

type AcceptChatRequest struct {
	AdminID uint `json:"admin_id" binding:"required"`
}

type SendMessageRequest struct {
	Content     string  `json:"content" binding:"required"`
	MessageType string  `json:"message_type" binding:"oneof=text image file"`
	MediaURL    *string `json:"media_url,omitempty"`
}

type GetChatHistoryRequest struct {
	SessionUID string `uri:"session_uid" binding:"required,uuid"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     int    `form:"offset" binding:"omitempty,min=0"`
}