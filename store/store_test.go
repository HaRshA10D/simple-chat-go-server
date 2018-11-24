package store

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

func TestCreateUser(t *testing.T) {

	db, err := gorm.Open("postgres", "user=shailendra dbname=simplechat sslmode=disable")
	if err != nil {
		t.Error("Not able to connect with database")
	}

	store := &SQLSupplier{
		simpleChatDatabase: db,
	}

	user := model.User{}

	store.simpleChatDatabase.CreateTable(&user)

	user1 := &model.User{
		Name: "Harry",
	}
	user2 := &model.User{
		Name: "Harry",
	}

	user1Token, user1Message := store.CreateUser(user1)
	user2Token, user2Message := store.CreateUser(user2)

	if user1Token == "" || user1Message != "login successful" {
		t.Error("new user should be add succesfully in database")
	}

	if user2Token == "" || user2Message != "login successful" {
		t.Error("existing user should not be add succesfully in database but should be login")
	}

	store.simpleChatDatabase.DropTable(&model.User{})
	store.simpleChatDatabase.Close()
}
