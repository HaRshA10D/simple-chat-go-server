package store

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

type SimpleChatStore interface {
	JoinGroup(user *model.User, group *model.Group) error
	CreateUser(user *model.User) error
	FindUserByToken(token string) (model.User, error)
	CreateGroup(group model.Group) (model.Group, error)
	UserGroups(user *model.User) ([]model.Group, error)
	InitDatabase() error
	DB() *gorm.DB
}

type sqlSupplier struct {
	db *gorm.DB
}

func (sqlSupplier *sqlSupplier) UserGroups(user *model.User) ([]model.Group, error) {
	resultGroups := []model.Group{}
	query := "SELECT groups.id, groups.name, groups.last_activity_at FROM user_groups INNER JOIN groups ON user_groups.group_id = groups.id  AND user_groups.user_id = ? ORDER BY groups.last_activity_at DESC"
	result := sqlSupplier.DB().Raw(query, user.ID).Scan(&resultGroups)
	if result.Error != nil {
		return nil, result.Error
	}
	return resultGroups, nil
}

func (sqlSupplier *sqlSupplier) JoinGroup(user *model.User, group *model.Group) error {
	resultGroup := &model.Group{}
	result := sqlSupplier.DB().Where("name = ?", group.Name).First(resultGroup)
	if result.RecordNotFound() {
		return errors.New("Group does not exist")
	}
	insertUserGroup := &model.UserGroup{
		UserID:  user.ID,
		GroupID: resultGroup.ID,
	}
	result = sqlSupplier.DB().Where("user_id = ? AND group_id = ?", insertUserGroup.UserID, insertUserGroup.GroupID).First(&model.UserGroup{})
	if !result.RecordNotFound() {
		return errors.New("You have already joined this group")
	}
	sqlSupplier.DB().Create(insertUserGroup)
	group.ID = insertUserGroup.GroupID
	return nil
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

func NewSimpleChatStore(config *model.Config) (SimpleChatStore, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=UTC", config.UserName, config.Password, config.DatabaseHost, config.DatabasePort, config.DBName)
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &sqlSupplier{
		db: db,
	}, nil
}

func (sqlSupplier *sqlSupplier) InitDatabase() error {
	user := &model.User{}
	group := &model.Group{}
	userGroup := &model.UserGroup{}
	groupMessage := &model.GroupMessage{}
	err := sqlSupplier.DB().CreateTable(user).Error
	if err != nil {
		return err
	}
	err = sqlSupplier.DB().CreateTable(group).Error
	if err != nil {
		return err
	}
	err = sqlSupplier.DB().CreateTable(userGroup).Error
	if err != nil {
		return err
	}
	err = sqlSupplier.DB().CreateTable(groupMessage).Error
	if err != nil {
		return err
	}
	return nil
}
