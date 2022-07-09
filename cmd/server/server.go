package main

import (
	"log"

	"github.com/nivanov045/silver-octo-train/cmd/server/api"
	"github.com/nivanov045/silver-octo-train/cmd/server/service"
	"github.com/nivanov045/silver-octo-train/cmd/server/storage"
)

func main() {
	storage := storage.New()

	serv := service.New(storage)

	myapi := api.New(serv)

	log.Fatalln(myapi.Run())
}
