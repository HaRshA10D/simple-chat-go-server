package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

func (api *API) InitRootRoutes() {
	api.Router.UserRouter.Handle("/login", api.ChatHandler(loginUser)).Methods("POST")
	api.Router.GroupRouter.Handle("", api.AuthRequiredChatHandler(createGroup)).Methods("POST")
	api.Router.GroupRouter.Handle("/{name}/join", api.AuthRequiredChatHandler(joinGroup)).Methods("POST")
	api.Router.GroupRouter.Handle("", api.AuthRequiredChatHandler(fetchUserGroups)).Methods("GET")
}

func loginUser(c *Context, w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	user := model.User{}
	json.NewDecoder(r.Body).Decode(&user)
	response := model.UserResponse{}
	response.Message = "Login Successful"
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

	if statusCode == 200 {
		c.Store.JoinGroup(c.User, &returnedGroup)
	}

	response["group_id"] = returnedGroup.ID
	response["group_name"] = returnedGroup.Name

	message := make(map[string]interface{})
	message["data"] = response
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

func joinGroup(c *Context, w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	vars := mux.Vars(r)
	group := &model.Group{
		Name: vars["name"],
	}

	response := make(map[string]interface{})
	err := c.Store.JoinGroup(c.User, group)
	if err != nil {
		if err.Error() == "Group does not exist" {
			statusCode = 404
		} else {
			statusCode = 409
		}
		response["message"] = err.Error()
	}

	response["group_id"] = group.ID
	response["group_name"] = group.Name

	message := make(map[string]interface{})
	message["data"] = response
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

func fetchUserGroups(c *Context, w http.ResponseWriter, r *http.Request) {
	statusCode := 200
	response := make(map[string]interface{})

	userGroups, err := c.Store.UserGroups(c.User)
	if err != nil {
		statusCode = 500
	}

	response["groups"] = userGroups
	message := make(map[string]interface{})
	message["data"] = response

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}
