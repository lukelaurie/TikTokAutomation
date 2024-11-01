package database

import (
	"net/http"
	"strings"

	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
)

func RetrieveUser(username string) (model.User, error) {
	var user model.User

	err := DB.First(&user, "username = ?", username).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func RegisterUser(username string, email string, password string, w http.ResponseWriter) bool {
	newUser := model.User{
		Username: username,
		Email: email,
		Password: password,
	}

	// insert the user into the database 
	err := DB.Create(&newUser).Error
	if err == nil {
		return false
	}

	if isUniqueViolation(err, "users_pkey") {
		http.Error(w, "username already exists", http.StatusConflict)
	} else if isUniqueViolation(err, "email") {
		http.Error(w, "email already exists", http.StatusConflict)
	} else {
		http.Error(nil, "database insert error", http.StatusInternalServerError)
	}
	return true
}

func isUniqueViolation(err error, constraint string) bool {
	return strings.Contains(err.Error(), "unique constraint") && 
		   strings.Contains(err.Error(), constraint)
}