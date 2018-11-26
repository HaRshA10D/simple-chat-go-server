package main

import (
	"flag"
	"log"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/model"

	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/api"
	"source.golabs.io/ops-tech-peeps/simple-chat-go-server/store"
)

func main() {
	migrate := flag.Bool("migrate", false, "To run migration")
	flag.Parse()
	config := model.NewConfig()
	store, err := store.NewSimpleChatStore(config)
	if err != nil {
		log.Fatalf("Can't connect to store : \n%s", err.Error())
	}
	if *migrate {
		err = store.InitDatabase()
		if err != nil {
			log.Fatalf("Could not perform migration due to: %s", err.Error())
		}
		if store.DB() != nil {
			store.DB().Close()
		}
		return
	}
	api := api.Init(store)
	api.StartServer(config.ServerPort, config.WriteTimeout, config.ReadTimeout, config.IdleTimeout)
	api.QuitSignalHandler(15)
}
