package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/lib/pq"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)

func RetrieveUser(username string) (model.User, error) {
	var user model.User
	query := `SELECT username, email, password FROM users WHERE username = $1`

	// search for the row 
	err := DB.QueryRow(query, username).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		// check if the error was from no user being found 
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("username invalid")
		}
		return user, err
	}
	return user, nil
}

func RegisterUser(username string, email string, password string, w http.ResponseWriter) bool {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := DB.Exec(query, username, email, password)
	if err != nil {
		// check if the error is a unique constraint violation
		pqErr, ok := err.(*pq.Error)
		if !ok || pqErr.Code != "23505" { // check for unique constraint violation
			http.Error(w, "database insert erro", http.StatusInternalServerError)
			return true
		}
		// check what violation occured
		if strings.Contains(pqErr.Message, "users_pkey") {
			http.Error(w, "username already exists", http.StatusConflict)
		} else if strings.Contains(pqErr.Message, "email") {
			http.Error(w, "email already exists", http.StatusConflict)
		}
		return true
	}

	return false
}