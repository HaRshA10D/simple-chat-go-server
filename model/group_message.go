package model

import (
	"time"
)

type GroupMessage struct {
	ID        string
	UserID    string 
	GroupID   string
	Message   string
	MessageSentAt time.Time
}
