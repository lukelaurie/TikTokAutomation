package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// global variable to manage the database connection
var DB *gorm.DB
var SqlDB *sql.DB

func InitDb() {
	dbPassword := os.Getenv("DB_PASSWORD")

	var err error
	// connect to the database
	connStr := fmt.Sprintf("host=localhost port=5432 user=postgres password=%s dbname=postgres sslmode=disable", dbPassword)

	// connect to the database with GORM 
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error opening connection to the database with GORM: %v", err))
	}

	// check the connection 
	SqlDB, err = DB.DB()
	if err != nil {
		panic(fmt.Errorf("error getting SQL DB from GORM: %v", err))
	}

	err = SqlDB.Ping()
	if err != nil {
		panic(fmt.Errorf("db ping error: %v", err))
	}

	// add or update the table with any changes made to the schema 
	updateTable()
}

func updateTable() error {
	var err error 

	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		return fmt.Errorf("error updating the User table in database: %v", err)
	}	

	err = DB.AutoMigrate(&model.Preference{})
	if err != nil {
		return fmt.Errorf("error updating the Preference table in database: %v", err)
	}	

	err = DB.AutoMigrate(&model.PreferenceTracker{})
	if err != nil {
		return fmt.Errorf("error updating the PreferenceTracker table in database: %v", err)
	}	

	return nil
}