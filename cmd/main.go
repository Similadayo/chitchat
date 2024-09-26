package main

import (
	"github/similadayo/chitchat/config"
	"log"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	config.MigrateDB(db)
}
