package model

type UserGroup struct {
	ID      int `gorm:"primary_key;SERIAL"`
	UserID  int
	GroupID int
}

func (userGroup *UserGroup) IsValid() bool {
	return !(userGroup.UserID == 0 || userGroup.GroupID == 0)
}
