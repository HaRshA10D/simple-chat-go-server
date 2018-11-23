package model

import (
	"regexp"
	"time"
)

type Group struct {
	ID             int
	Name           string
	LastActivityAt time.Time
}

func (group *Group) IsValid() bool {
	alphaNumericRegex, _ := regexp.Compile("^[a-zA-Z0-9_]+$")
	return alphaNumericRegex.MatchString(group.Name)
}
