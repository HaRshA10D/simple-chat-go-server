package store

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

type SimpleChatStore interface {
	CreateUser(user *model.User)
}

type SQLSupplier struct {
	simpleChatDatabase *gorm.DB
}

func (sqlSupplier *SQLSupplier) CreateUser(currentUser *model.User) (string, string) {
	user := model.User{}
	result := sqlSupplier.simpleChatDatabase.First(&user, "name = ?", currentUser.Name)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			currentUser.Token = uuid.NewV4().String()
			err := sqlSupplier.simpleChatDatabase.Create(currentUser).Error
			if err != nil {
				return "", "not login successfull"
			}
			return currentUser.Token, "login successful"
		}
		return "", "not login successfull"
	}
	return user.Token, "login successful"
}
