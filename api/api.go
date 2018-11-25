package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/store"
)

type API struct {
	Store  store.SimpleChatStore
	Router *Routes
	Server *http.Server
}

type Routes struct {
	RootRouter  *mux.Router
	UserRouter  *mux.Router
	GroupRouter *mux.Router
}

func Init(store store.SimpleChatStore) *API {
	api := &API{
		Store:  store,
		Router: &Routes{},
	}
	api.Router.RootRouter = mux.NewRouter()
	api.Router.UserRouter = api.Router.RootRouter.PathPrefix("/users").Subrouter()
	api.Router.GroupRouter = api.Router.RootRouter.PathPrefix("/groups").Subrouter()
	api.InitRootRoutes()
	return api
}

func (api *API) StartServer(port, WriteTimeout, ReadTimeout, IdleTimeout int) {
	var handler http.Handler = api.Router.RootRouter
	api.Server = &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Duration(WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(IdleTimeout) * time.Second,
		Handler:      handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handler),
	}
	log.Printf("Listening to port : %d \n", port)
	go api.Server.ListenAndServe()
}

func (api *API) QuitSignalHandler(ShutdownTimeout int) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	for {
		select {
		case sig := <-c:
			log.Printf("Got %s signal. Aborting...\n", sig)
			api.StopHttpServer(ShutdownTimeout)
			api.Store.DB().Close()
			return
		}
	}
}

func (api *API) StopHttpServer(ShutdownTimeout int) {
	if api.Server != nil {
		log.Print("Stopping Server! Bye.")
		defer func() {
			api.Server = nil
		}()
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ShutdownTimeout)*time.Second*10)
		defer cancel()
		api.Server.Shutdown(ctx)
	}
}
