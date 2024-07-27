package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
	database "github.com/lukelaurie/TikTokAutomation/backend/internal/database"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// execute the query to add a new user
	username := "test username2"
	email := "testemail@google.com2"
	password := "a password"

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, username, email, password)
	if err != nil {
		panic(fmt.Errorf("database insert error: %v", err))
	}

	json.NewEncoder(w).Encode("user regiserd")
}