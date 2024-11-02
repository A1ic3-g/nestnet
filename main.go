package main

import (
	"log"
	"nestnet/internal/database"
)

func main() {
	//service.Start()

	log.Println(database.GetPosts())
}
