package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

func (api *API) InitRootRoutes() {
	api.Router.UserRouter.Handle("/login", api.ChatHandler(loginUser)).Methods("POST")
	api.Router.GroupRouter.Handle("", api.AuthRequiredChatHandler(createGroup)).Methods("POST")
	api.Router.GroupRouter.Handle("/{id}/messages", api.AuthRequiredChatHandler(sendMessageToGroup)).Methods("POST")
}

type MessageRequest struct {
	Text              string
	Message_sent_time string
}

func sendMessageToGroup(c *Context, w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	var messageRequest MessageRequest
	vars := mux.Vars(r)
	id := vars["id"]
	json.NewDecoder(r.Body).Decode(&messageRequest)
	w.Header().Set("Content-type", "application/json")
	message := make(map[string]interface{})
	response := make(map[string]interface{})

	epochMillis, err := strconv.ParseInt(messageRequest.Message_sent_time, 10, 64)
	if err != nil {
		response["message"] = "Bad Request for Message Sent time"
		message["data"] = response
		statusCode = 400
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(message)
		return
	}

	group, err := c.Store.FindGroupByID(id)
	if err != nil {
		response["message"] = "Group does not exist"
		message["data"] = response
		statusCode = 404
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(message)
		return
	}

	err = c.Store.SendMessage(c.User, &group, messageRequest.Text, epochMillis)

	if err != nil {
		response["message"] = "Internal Server Error"
		message["data"] = response
		statusCode = 500
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(message)
		return
	}

	response["message"] = "Successfully sent message"
	message["data"] = response
	statusCode = 200
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)

}

func loginUser(c *Context, w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	user := model.User{}
	json.NewDecoder(r.Body).Decode(&user)
	response := model.UserResponse{}
	err := c.Store.CreateUser(&user)
	if err != nil {
		response.Message = "Not able to login"
		statusCode = 400
	}
	response.Token = user.Token
	message := make(map[string]interface{})
	message["data"] = response
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

func createGroup(c *Context, w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	group := model.Group{}
	json.NewDecoder(r.Body).Decode(&group)
	response := make(map[string]interface{})
	returnedGroup, err := c.Store.CreateGroup(group)
	if err != nil {
		response["message"] = "Not able to create group"
		statusCode = 500
	}
	response["group_id"] = returnedGroup.ID
	response["group_name"] = returnedGroup.Name

	message := make(map[string]interface{})
	message["data"] = response
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}
