package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	// TODO generate cookie 
	json.NewEncoder(w).Encode("login successful")
}
