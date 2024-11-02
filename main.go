package main

import (
	"log"
	"nestnet/internal/database"
	"nestnet/internal/service"
)

func main() {
	service.Start()

	log.Println(database.GetPosts())
}
