package main

import (
	"log"
	"os"

	"github.com/thisisheymde/URL-shortener/backend/api"
	"github.com/thisisheymde/URL-shortener/backend/storage"
)

func main() {
	store, err := storage.StartRedis(os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(":8081", store)
	server.Run()
}
