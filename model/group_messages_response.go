package model

type GroupMessageResponse struct {
	Message       string `json:"message"`
	MessageSentAt int64  `json:"message_sent_time"`
	UserName      string `json:"user_name"`
}
