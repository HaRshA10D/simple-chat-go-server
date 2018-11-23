package model

import (
	"time"
)

type GroupMessage struct {
	ID        int
	UserID    int 
	GroupID   int
	Message   string
	MessageSentAt time.Time
}
