package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/mocks"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

func TestLoginUser(t *testing.T) {

	simpleChatStore := &mocks.SimpleChatStore{}

	api := &API{}
	api.Store = simpleChatStore

	//FIXME: return token from create user along with error
	simpleChatStore.On("CreateUser", mock.MatchedBy(func(input *model.User) bool {
		input.Token = "1234"
		return true
	})).Return(nil).Once()

	userName := []byte(`{"name":"amit"}`)
	url := "localhost:3000/users/login"
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(userName))
	handler := api.ChatHandler(loginUser)
	handler.ServeHTTP(rr, req)
	var getResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &getResponse)
	dataByte, _ := json.Marshal(getResponse["data"])
	returnUserResponse := &model.UserResponse{}
	json.Unmarshal(dataByte, &returnUserResponse)
	assert.Equal(t, http.StatusOK, rr.Code, "Should get code 200 when user is created")
	assert.Equal(t, returnUserResponse.Token, "1234", "should get valid token")

	simpleChatStore.On("CreateUser", mock.Anything).Return(errors.New("Not able to create user")).Once()
	rr = httptest.NewRecorder()

	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(userName))
	handler = api.ChatHandler(loginUser)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Should return 400 bad request when create user failed")
}

func TestCreateGroup(t *testing.T) {

	simpleChatStore := &mocks.SimpleChatStore{}

	api := &API{}
	api.Store = simpleChatStore

	group := model.Group{
		ID:   10,
		Name: "group",
	}

	existinguser := model.User{
		ID:    10,
		Token: "12345",
		Name:  "Harry",
	}

	//FIXME: return token from create user along with error
	simpleChatStore.On("CreateGroup", mock.Anything).Return(group, nil).Once()
	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existinguser, nil).Once()

	groupName := []byte(`{"name":"group"}`)
	url := "localhost:3000/groups"
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(groupName))
	req.Header.Add("Auth-Token", "token "+existinguser.Token)
	handler := api.AuthRequiredChatHandler(createGroup)
	handler.ServeHTTP(rr, req)
	var returnResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &returnResponse)
	dataByte, _ := json.Marshal(returnResponse["data"])
	var parsedReturnResponse map[string]interface{}
	json.Unmarshal(dataByte, &parsedReturnResponse)
	dataByte1, _ := json.Marshal(parsedReturnResponse["group_id"])

	assert.Equal(t, http.StatusOK, rr.Code, "Should get code 200 when user is created")
	assert.Equal(t, fmt.Sprintf("%v", group.ID), string(dataByte1), "New group should be created")
}
