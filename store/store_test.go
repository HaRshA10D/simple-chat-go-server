package store

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

func TestCreateUser(t *testing.T) {

	user := model.User{}

	store, storeErr := setUp(user)
	if storeErr != nil {
		t.Fatalf("Not able to create store: %s", storeErr.Error())
	}
	defer tearDown(store.DB(), user)

	user1 := &model.User{
		Name: "Harry",
	}
	user2 := &model.User{
		Name: "Harry",
	}

	err1 := store.CreateUser(user1)
	err2 := store.CreateUser(user2)

	if err1 != nil || user1.Token == "" {
		t.Error("new user should be add succesfully in database")
	}

	if err2 != nil || user2.Token == "" {
		t.Error("existing user should not be add succesfully in database but should be login")
	}
}

func TestFindUser(t *testing.T) {

	user := model.User{}

	store, storeErr := setUp(user)
	if storeErr != nil {
		t.Fatalf("Not able to create store: %s", storeErr.Error())
	}
	defer tearDown(store.DB(), user)

	token := "123456789"

	_, findErr := store.FindUserByToken(token)

	if findErr == nil {
		t.Errorf("Should return error when token is not found")
	}

	user1 := &model.User{
		Name: "Harry",
	}

	store.CreateUser(user1)

	returnedUser, findErr := store.FindUserByToken(user1.Token)

	if findErr != nil {
		t.Errorf("should return user when we find user with valid token %s", findErr.Error())
	}

	if user1.ID != returnedUser.ID {
		t.Error("Invalid user returned")
	}

}

func TestCreateGroup(t *testing.T) {

	user := model.User{}
	group := model.Group{}

	store, storeErr := setUp(user, group)
	if storeErr != nil {
		t.Fatalf("Not able to create store: %s", storeErr.Error())
	}
	defer tearDown(store.DB(), user, group)

	token := "123456789"

	_, findErr := store.FindUserByToken(token)

	if findErr == nil {
		t.Errorf("Should return error when token is not found")
	}

	user1 := &model.User{
		Name: "Harry",
	}

	store.CreateUser(user1)

	group1 := model.Group{
		Name: "GroupName",
	}
	returnedGroup, createGroupError := store.CreateGroup(group1)
	if createGroupError != nil {
		t.Errorf("Create Group Should Create Group %s", createGroupError.Error())
	}
	if returnedGroup.Name != group1.Name && returnedGroup.ID == 0 {
		t.Error("Create Group Should Create Group")
	}
}

func TestJoinGroup(t *testing.T) {

	group := &model.Group{}
	user := &model.User{}
	userGroup := &model.UserGroup{}

	store, storeErr := setUp(user, group, userGroup)
	if storeErr != nil {
		t.Fatalf("Not able to create store: %s", storeErr.Error())
	}
	defer tearDown(store.DB(), user, group, userGroup)

	testUser := &model.User{
		ID:   5,
		Name: "Pinkman",
	}
	testGroup := &model.Group{
		Name: "Mexico",
	}

	err1 := store.JoinGroup(testUser, testGroup)
	assert.Error(t, err1, "Should not allow to join a non existent group")

	insertGroup := &model.Group{
		ID:   20,
		Name: "Mexico",
	}
	store.DB().Create(insertGroup)

	_ = store.JoinGroup(testUser, testGroup)

	findUserGroup := store.DB().Where("user_id = ? AND group_id = ?", testUser.ID, insertGroup.ID).First(&model.UserGroup{})
	assert.False(t, findUserGroup.RecordNotFound(), "Should join a user in to given group")

	assert.Equal(t, insertGroup.ID, testGroup.ID, "Should aso find the group id")

	err2 := store.JoinGroup(testUser, testGroup)
	assert.Error(t, err2, "Shouldn't allow a user to join group more than once")
}

func setUp(tables ...interface{}) (SimpleChatStore, error) {
	config := model.NewConfig()

	store, storeErr := NewSimpleChatStore(config)

	if storeErr != nil {
		return nil, storeErr
	}
	for _, table := range tables {
		store.DB().CreateTable(table)
	}
	return store, nil
}

func tearDown(db *gorm.DB, tables ...interface{}) {
	for _, table := range tables {
		db.DropTable(table)
	}
	db.Close()
}
