package main

import (
	"log"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/api"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/store"
)

func main() {
	config := model.NewConfig()
	store, err := store.NewSimpleChatStore(config)
	if err != nil {
		log.Fatalf("Can't connect to store : \n%s", err.Error())
	}

	api := api.Init(store)
	api.StartServer(config.ServerPort, config.WriteTimeout, config.ReadTimeout, config.IdleTimeout)
	api.QuitSignalHandler(15)
}
