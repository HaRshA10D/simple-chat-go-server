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
		Token: "TOKENVALUE",
		Name:  "Harry",
	}

	//FIXME: return token from create user along with error
	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, errors.New("No Auth")).Once()

	groupName := []byte(`{"name":"group"}`)
	url := "localhost:3000/groups"
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(groupName))
	handler := api.AuthRequiredChatHandler(createGroup)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code, "Should get code 401 when user is not authorised")

	simpleChatStore.On("CreateGroup", mock.Anything).Return(group, nil).Once()
	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()
	simpleChatStore.On("JoinGroup", mock.Anything, mock.Anything).Return(nil).Once()
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(groupName))
	req.Header.Set("Auth-Token", "TOKENVALUE")
	handler = api.AuthRequiredChatHandler(createGroup)

	handler.ServeHTTP(rr, req)
	var returnResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &returnResponse)
	dataByte, _ := json.Marshal(returnResponse["data"])
	var parsedReturnResponse map[string]interface{}
	json.Unmarshal(dataByte, &parsedReturnResponse)
	dataByte1, _ := json.Marshal(parsedReturnResponse["group_id"])

	assert.Equal(t, http.StatusOK, rr.Code, "Should get code 200 when user is created")
	assert.Equal(t, fmt.Sprintf("%v", group.ID), string(dataByte1), "New group should be created")

	simpleChatStore.On("CreateGroup", mock.Anything).Return(group, errors.New("Internal Database error")).Once()
	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(groupName))
	req.Header.Set("Auth-Token", "TOKENVALUE")
	handler = api.AuthRequiredChatHandler(createGroup)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code, "Should get code 500 when group is created")
}

func TestSendMessage(t *testing.T) {
	simpleChatStore := &mocks.SimpleChatStore{}

	api := &API{}
	api.Store = simpleChatStore

	group := model.Group{
		ID:   10,
		Name: "group",
	}

	existingUser := model.User{
		ID:    10,
		Token: "TOKENVALUE",
		Name:  "Harry",
	}

	requestData := []byte(`
	{
		"text":"Messsage Text",
		"message_sent_time":"1543217006000"
	}
	`)

	requestDataBadTime := []byte(`
	{
		"text":"Messsage Text",
		"message_sent_time":"Abhinav"
	}
	`)

	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, errors.New("No Auth")).Once()
	url := "localhost:3000/groups/10/messages/"
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	req = mux.SetURLVars(req, map[string]string{"id": "10"})
	handler := api.AuthRequiredChatHandler(sendMessageToGroup)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code, "Should get code 401 when user is not authorised")

	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()
	url = "localhost:3000/groups/10/messages"
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(requestDataBadTime))
	req = mux.SetURLVars(req, map[string]string{"id": "10"})
	req.Header.Set("Auth-Token", "TOKENVALUE")
	handler = api.AuthRequiredChatHandler(sendMessageToGroup)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 400, rr.Code, "Should get code 400 when message sent time is bad")

	simpleChatStore.On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	simpleChatStore.On("FindGroupByID", mock.Anything).Return(group, nil).Once()
	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()
	url = "localhost:3000/groups/10/messages"
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	req = mux.SetURLVars(req, map[string]string{"id": "10"})
	req.Header.Set("Auth-Token", "TOKENVALUE")
	handler = api.AuthRequiredChatHandler(sendMessageToGroup)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Should get code 200 when message is sent")

	simpleChatStore.On("FindGroupByID", mock.Anything).Return(group, errors.New("group does not exist")).Once()
	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()
	url = "localhost:3000/groups/1/messages"
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req.Header.Set("Auth-Token", "TOKENVALUE")
	handler = api.AuthRequiredChatHandler(sendMessageToGroup)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 404, rr.Code, "Should get code 404 when group is not present")

	simpleChatStore.On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Internal Database Error")).Once()
	simpleChatStore.On("FindGroupByID", mock.Anything).Return(group, nil).Once()
	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Once()
	url = "localhost:3000/groups/10/messages"
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	req = mux.SetURLVars(req, map[string]string{"id": "10"})
	req.Header.Set("Auth-Token", "TOKENVALUE")
	handler = api.AuthRequiredChatHandler(sendMessageToGroup)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 500, rr.Code, "Should get code 500 when internal db error occurs")
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

func TestFetchUserGroups(t *testing.T) {

	simpleChatStore := &mocks.SimpleChatStore{}
	api := &API{}
	api.Store = simpleChatStore

	existingUser := model.User{
		ID:   10,
		Name: "Harry",
	}
	group1 := model.Group{
		ID:   1,
		Name: "group1",
	}
	group2 := model.Group{
		ID:   2,
		Name: "group2",
	}
	group3 := model.Group{
		ID:   3,
		Name: "group3",
	}
	group4 := model.Group{
		ID:   4,
		Name: "group4",
	}

	userGroups := []model.Group{group1, group2, group3, group4}

	simpleChatStore.On("FindUserByToken", mock.Anything).Return(existingUser, nil).Twice()
	simpleChatStore.On("UserGroups", mock.Anything).Return(userGroups, nil).Once()

	url := "localhoset:9090/groups"
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Auth-Token", "12345")
	handler := api.AuthRequiredChatHandler(fetchUserGroups)
	handler.ServeHTTP(rr, req)

	var returnResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &returnResponse)
	dataByte, _ := json.Marshal(returnResponse["data"])
	var parsedReturnResponse map[string]interface{}
	json.Unmarshal(dataByte, &parsedReturnResponse)

	responseGroupByte, _ := json.Marshal(parsedReturnResponse["groups"])
	responseGroup := []model.Group{}
	json.Unmarshal(responseGroupByte, &responseGroup)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	assert.Equal(t, len(userGroups), len(responseGroup), "Should return 5 group details")

	simpleChatStore.On("UserGroups", mock.Anything).Return([]model.Group{}, errors.New("Internal error")).Once()

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Auth-Token", "12345")
	handler = api.AuthRequiredChatHandler(fetchUserGroups)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should throw internal error")
}
