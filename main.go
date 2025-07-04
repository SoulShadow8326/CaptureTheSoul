package main

import(
	"log"
	"net/http"

	"CaptureTheSoul/database"
	"CaptureTheSoul/routes"
)

func main(){
	database.InitDB("data.db")

	routes.SetupRoutes()

	log.Println("http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		log.Fatal("haha your server broken:", err)
	}
}