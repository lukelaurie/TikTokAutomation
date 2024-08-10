package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/database"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/utils"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var reqBody model.LoginRequest

	// Decode the body of the request into the struct
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		panic(fmt.Errorf("request decode error: %v", err))
	}

	// search for the user in the database
	user, err := database.RetrieveUser(reqBody.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// check that the user passed in the correct password
	if !utils.CheckPasswordhash(reqBody.Password, user.Password) {
		http.Error(w, "password invalid", http.StatusConflict)
		return
	}

	// create the JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // set how long token lasts
	})

	// sign the token with a secret key
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	// encodes both the username and exp into tokenString
	tokenString, err := jwtToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		http.Error(w, "unable to generate the jwt token", http.StatusInternalServerError)
		return
	}

	// place the token in a cookie in the request
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: time.Now().Add(time.Hour * 72),
		HttpOnly: true,
		Secure: true,
		Path: "/",
	})

	json.NewEncoder(w).Encode("login successful")
}