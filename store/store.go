package store

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

type SimpleChatStore interface {
	CreateUser(user *model.User) error
	FindUserByToken(token string) (model.User, error)
	CreateGroup(group model.Group) (model.Group, error)
	DB() *gorm.DB
}

type sqlSupplier struct {
	db *gorm.DB
}

func (sqlSupplier *sqlSupplier) CreateGroup(group model.Group) (model.Group, error) {
	result := sqlSupplier.DB().First(&model.Group{}, "name = ?", group.Name)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			err := sqlSupplier.DB().Create(&group).Error
			if err != nil {
				return model.Group{}, err
			}
			return group, nil
		}
		return model.Group{}, err
	}
	return model.Group{}, errors.New("Group Name already exists")
}

func (sqlSupplier *sqlSupplier) FindUserByToken(token string) (model.User, error) {
	user := model.User{}
	result := sqlSupplier.DB().First(&user, "token = ?", token)
	if err := result.Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (sqlSupplier *sqlSupplier) CreateUser(currentUser *model.User) error {
	user := model.User{}
	result := sqlSupplier.DB().First(&user, "name = ?", currentUser.Name)
	if err := result.Error; err != nil {
		if result.RecordNotFound() {
			currentUser.Token = uuid.NewV4().String()
			err := sqlSupplier.DB().Create(currentUser).Error
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	currentUser.Token = user.Token
	currentUser.ID = user.ID
	return nil
}

func (sqlSupplier *sqlSupplier) DB() *gorm.DB {
	return sqlSupplier.db
}

func NewSimpleChatStore(dbHost, dbPort, dbName, dbUser, dbPass string) (SimpleChatStore, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=UTC", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &sqlSupplier{
		db: db,
	}, nil
}
