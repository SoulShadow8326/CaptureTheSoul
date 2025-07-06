package main

import (
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"CaptureTheSoul/database"
	"CaptureTheSoul/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system env")
	}

	database.InitDB("data.db")
	routes.SetupRoutes()

	log.Println("http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("your server broke:", err)
	}
}
