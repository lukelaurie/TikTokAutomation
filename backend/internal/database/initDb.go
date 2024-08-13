package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// global variable to manage the database connection
var DB *sql.DB

func InitDb() {
	dbPassword := os.Getenv("DB_PASSWORD")

	var err error
	// connect to the database
	connStr := fmt.Sprintf("host=localhost port=5432 user=postgres password=%s dbname=postgres sslmode=disable", dbPassword)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Errorf("connection open error: %v", err))
	}
	
	// check the connection 
	err = DB.Ping()
	if err != nil {
		panic(fmt.Errorf("db ping error: %v", err))
	}
}