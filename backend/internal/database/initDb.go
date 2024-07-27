package database 

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// global variable to manage the database connection 
var DB *sql.DB

func InitDb() {
	var err error
	// connect to the database
	connStr := "host=localhost port=5432 user=postgres password=Sabinois@1225 dbname=postgres sslmode=disable"

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