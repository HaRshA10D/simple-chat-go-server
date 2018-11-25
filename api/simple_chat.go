package api

import (
	"encoding/json"
	"net/http"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

func (api *API) InitRootRoutes() {
	api.Router.UserRouter.Handle("/login", api.ChatHandler(loginUser)).Methods("POST")
	api.Router.GroupRouter.Handle("/", api.AuthRequiredChatHandler(createGroup)).Methods("POST")
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
		response["message"] = "Not able to craete group"
		statusCode = 400
	}
	response["group_id"] = returnedGroup.ID
	response["group_name"] = returnedGroup.Name

	message := make(map[string]interface{})
	message["data"] = response
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}
