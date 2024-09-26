package main

import (
	"github/similadayo/chitchat/config"
	"github/similadayo/chitchat/routes"
	"log"
	"net/http"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	config.MigrateDB(db)

	log.Println("Successfully connected to the database")

	r := routes.InitRoutes()

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start the server: %v", err)
	}
}
