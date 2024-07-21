package main

import (
	"log"
	"net/http"

	route "github.com/lukelaurie/TikTokAutomation/backend/internal/routes"
)

func main() {
	router := route.InitializeRoutes()
	log.Println("Server started on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}