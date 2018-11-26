package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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

	existingUser := model.User{
		ID:    10,
		Token: "12345",
		Name:  "Harry",
	}

	//FIXME: return token from create user along with error
	simpleChatStore.On("CreateGroup", mock.Anything).Return(group, nil).Once()
	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()

	groupName := []byte(`{"name":"group"}`)
	url := "localhost:3000/groups"
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(groupName))

	handler := api.ChatHandler(createGroup)
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

func TestJoinGroup(t *testing.T) {

	simpleChatStore := &mocks.SimpleChatStore{}
	api := &API{}
	api.Store = simpleChatStore

	group := &model.Group{
		ID:   10,
		Name: "fun",
	}
	existingUser := model.User{
		ID:   10,
		Name: "Harry",
	}
	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()
	simpleChatStore.On("JoinGroup", mock.Anything, mock.MatchedBy(func(gr *model.Group) bool {
		gr.ID = group.ID
		return true
	})).Return(nil).Once()

	url := "localhoset:9090/groups/fun/join"
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, nil)
	req = mux.SetURLVars(req, map[string]string{"name": "fun"})
	req.Header.Set("Auth-Token", "12345")
	handler := api.AuthRequiredChatHandler(joinGroup)
	handler.ServeHTTP(rr, req)

	var returnResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &returnResponse)
	dataByte, _ := json.Marshal(returnResponse["data"])
	var parsedReturnResponse map[string]interface{}
	json.Unmarshal(dataByte, &parsedReturnResponse)

	groupIDFloat, ok := parsedReturnResponse["group_id"].(float64)
	assert.True(t, ok, "Expected group ID to be Int")
	groupID := int(groupIDFloat)
	groupName, ok := parsedReturnResponse["group_name"].(string)

	assert.True(t, ok, "Expected group Name to be string")
	assert.Equal(t, http.StatusOK, rr.Code, "Should get code 200 when user joined the group")
	assert.Equal(t, group.ID, groupID, "Expected to return group ID as 10")
	assert.Equal(t, group.Name, groupName, "Expected group name to be same as sent")

	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()
	simpleChatStore.On("JoinGroup", mock.Anything, mock.Anything).Return(errors.New("Group does not exist")).Once()

	rr = httptest.NewRecorder()
	handler1 := api.AuthRequiredChatHandler(joinGroup)
	handler1.ServeHTTP(rr, req)

	json.Unmarshal(rr.Body.Bytes(), &returnResponse)
	dataByte, _ = json.Marshal(returnResponse["data"])
	json.Unmarshal(dataByte, &parsedReturnResponse)

	message, _ := parsedReturnResponse["message"].(string)
	assert.Equal(t, http.StatusNotFound, rr.Code, "Should get code 200 when user joined the group")
	assert.Equal(t, "Group does not exist", message, "Expected a group dosen't exist error")

	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()
	simpleChatStore.On("JoinGroup", mock.Anything, mock.Anything).Return(errors.New("Already a member")).Once()

	rr = httptest.NewRecorder()
	handler2 := api.AuthRequiredChatHandler(joinGroup)
	handler2.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code, "Should not join a user in group user is already part of")
}
