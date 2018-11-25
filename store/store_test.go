package store

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

func TestCreateUser(t *testing.T) {

	store, err := NewSimpleChatStore("localhost", "5432", "simplechat", "shailendra", "")

	if err != nil {
		t.Fatalf("not able to create store %s", err.Error())
	}

	user := model.User{}

	store.DB().CreateTable(&user)

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

	store.DB().DropTable(&model.User{})
	store.DB().Close()
}

func TestFindUser(t *testing.T) {
	store, err := NewSimpleChatStore("localhost", "5432", "simplechat", "shailendra", "")

	if err != nil {
		t.Fatalf("not able to create store %s", err.Error())
	}

	user := model.User{}

	store.DB().CreateTable(&user)

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
	store.DB().DropTable(&model.User{})
	store.DB().Close()

}

func TestCreateGroup(t *testing.T) {
	store, err := NewSimpleChatStore("localhost", "5432", "simplechat", "shailendra", "")

	if err != nil {
		t.Fatalf("not able to create store %s", err.Error())
	}

	user := model.User{}

	store.DB().CreateTable(&user)

	token := "123456789"

	_, findErr := store.FindUserByToken(token)

	if findErr == nil {
		t.Errorf("Should return error when token is not found")
	}

	user1 := &model.User{
		Name: "Harry",
	}

	store.CreateUser(user1)

	group := model.Group{}
	store.DB().CreateTable(group)

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
	store.DB().DropTable(&model.User{})
	store.DB().DropTable(&model.Group{})
	store.DB().Close()
}
