package api

import (
	"net/http"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/store"
)

type Context struct {
	Store store.SimpleChatStore
	User  *model.User
}

type Handler struct {
	Store                  store.SimpleChatStore
	HandlerFunc            func(*Context, http.ResponseWriter, *http.Request)
	AuthenticationRequired bool
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{}
	c.Store = h.Store

	if h.AuthenticationRequired {
		token := r.Header.Get("Auth-Token")
		user, err := h.Store.FindUserByToken(token)
		if err != nil {
			w.WriteHeader(401)
			// TODO: convert to JSON
			w.Write([]byte("Not authorised"))
			return
		}
		c.User = &user
	}
	h.HandlerFunc(c, w, r)
}

func (a *API) ChatHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &Handler{
		Store:                  a.Store,
		HandlerFunc:            h,
		AuthenticationRequired: false,
	}
}

func (a *API) AuthRequiredChatHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &Handler{
		Store:                  a.Store,
		HandlerFunc:            h,
		AuthenticationRequired: true,
	}
}
