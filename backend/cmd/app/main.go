package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/database"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/routes"
)

func init() {
	// load in the .env file 
	err := godotenv.Load("./cmd/app/.env")
	if err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}
}

func main() {
	// check command line to see if should run in test mode 
	var isTestMode bool
	flag.BoolVar(&isTestMode, "test", false, "Run server in test mode") // set to true if test command passed
	flag.Parse()

	// initialize the database connection
	database.InitDb()
	defer database.SqlDB.Close()
	
	// set up the router
	router := route.InitializeRoutes(isTestMode)
	log.Println("Server started on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}