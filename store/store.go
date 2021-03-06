package store

import (
	"errors"
	"fmt"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

type SimpleChatStore interface {
	JoinGroup(user *model.User, group *model.Group) error
	CreateGroup(group model.Group) (model.Group, error)
	CreateUser(user *model.User) error
	FindUserByToken(token string) (model.User, error)
	FindGroupByID(ID string) (model.Group, error)
	SendMessage(user *model.User, group *model.Group, messageText string, sendMessageTime int64) error
	UserGroups(user *model.User) ([]model.Group, error)
	GroupMessages(user *model.User, groupID int) ([]model.GroupMessage, error)
	InitDatabase() error
	DB() *gorm.DB
}

type sqlSupplier struct {
	db *gorm.DB
}

func (sqlSupplier *sqlSupplier) SendMessage(user *model.User, group *model.Group, messageText string, sendMessageTime int64) error {

	returnUser := sqlSupplier.DB().First(&model.User{}, "name = ?", user.Name)

	if err := returnUser.Error; err != nil {
		return errors.New("User does not exist for send message")
	}

	returnGroup := sqlSupplier.DB().First(&model.Group{}, "ID = ?", group.ID)

	if err := returnGroup.Error; err != nil {
		return errors.New("Group does not exist for send message")
	}

	groupMessage := &model.GroupMessage{
		UserID:        user.ID,
		GroupID:       group.ID,
		Message:       messageText,
		MessageSentAt: sendMessageTime,
		UserName:      user.Name,
	}

	err := sqlSupplier.DB().Create(groupMessage).Error
	if err != nil {
		return err
	}
	return nil
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

func (sqlSupplier *sqlSupplier) FindGroupByID(ID string) (model.Group, error) {
	group := model.Group{}
	result := sqlSupplier.DB().First(&group, "ID = ?", ID)
	if err := result.Error; err != nil {
		return model.Group{}, err
	}
	return group, nil
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

func (sqlSupplier *sqlSupplier) GroupMessages(user *model.User, groupID int) ([]model.GroupMessage, error) {
	returnUser := sqlSupplier.DB().First(&model.User{}, "name = ?", user.Name)
	if err := returnUser.Error; err != nil {
		return nil, errors.New("User does not exist for requested group messages")
	}

	returnGroup := sqlSupplier.DB().First(&model.Group{}, "ID = ?", groupID)
	if err := returnGroup.Error; err != nil {
		return nil, errors.New("Group does not exist for requested group messages")
	}

	returnUserGroup := sqlSupplier.DB().Where(&model.UserGroup{}, "UserID = ? AND GroupID = ?", user.ID, groupID).First(&model.UserGroup{})
	if err := returnUserGroup.Error; err != nil {
		return nil, errors.New("UserGroup does not exist for requested group messages")
	}

	userMessages := []model.GroupMessage{}
	userMessagesFromDb := sqlSupplier.DB().Raw("SELECT user_name, message, message_sent_at FROM group_messages WHERE group_id = ? order by message_sent_at DESC limit 10", groupID).Scan(&userMessages)
	if err := userMessagesFromDb.Error; err != nil {
		return nil, errors.New("Internal DB Error")
	}
	return userMessages, nil
}

func (sqlSupplier *sqlSupplier) DB() *gorm.DB {
	return sqlSupplier.db
}

func NewSimpleChatStore(config *model.Config) (SimpleChatStore, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=UTC", *config.DatabaseSettings.UserName, *config.DatabaseSettings.Password, *config.DatabaseSettings.DatabaseHost, *config.DatabaseSettings.DatabasePort, *config.DatabaseSettings.DBName)
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
