package api

import (
	"encoding/json"
	"net/http"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"
)

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
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}
