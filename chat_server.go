package main

import (
	"log"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/api"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/store"
)

func main() {
	store, err := store.NewSimpleChatStore("localhost", "5432", "simplechat", "shailendra", "")
	if err != nil {
		log.Fatalf("Can't connect to store : \n%s", err.Error())
	}
	api := api.Init(store)
	api.StartServer(9090, 300, 300, 300)
	api.QuitSignalHandler(15)
}
