package model

import (
	"regexp"
	"time"
)

type Group struct {
	ID             int       `gorm:"primary_key;SERIAL" json:"id"`
	Name           string    `json:"name"`
	LastActivityAt time.Time `json:"last_activity_at"`
}

func (group *Group) IsValid() bool {
	alphaNumericRegex, _ := regexp.Compile("^[a-zA-Z0-9_]+$")
	return alphaNumericRegex.MatchString(group.Name)
}
