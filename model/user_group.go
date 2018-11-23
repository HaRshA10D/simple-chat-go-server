package model

type UserGroup struct {
	ID      int
	UserID  int
	GroupID int
}

func (userGroup *UserGroup) IsValid() bool {
	return !(userGroup.UserID == 0 || userGroup.GroupID == 0)
}
