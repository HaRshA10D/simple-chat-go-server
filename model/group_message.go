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

func (groupMessage *GroupMessage) IsValid() bool {
	return !(groupMessage.UserID == 0 || groupMessage.GroupID == 0 || groupMessage.Message == "")
}
