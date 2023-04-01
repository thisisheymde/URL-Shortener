package main

import (
	"log"

	"github.com/thisisheymde/URL-shortener/backend/api"
	"github.com/thisisheymde/URL-shortener/backend/storage"
)

func main() {
	store, err := storage.StartRedis("containers-us-west-33.railway.app:7772", "pOK9WhpY0SZ5TLh8T6ui")

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(":8081", store)
	server.Run()
}
