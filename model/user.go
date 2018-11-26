package model

import (
	"regexp"
)

type User struct {
	ID    int `gorm:"primary_key;SERIAL"`
	Name  string
	Token string
}

func (user *User) IsValid() bool {
	alphaNumericRegex, _ := regexp.Compile("^[a-zA-Z0-9]+$")
	return alphaNumericRegex.MatchString(user.Name)
}
