package handler

import (
	"encoding/json"
	"fmt"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/database"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/utils"
	"net/http"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var reqBody model.User

	// Decode the body of the request into the struct
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		utils.LogAndAddServerError(fmt.Errorf("request decode error: %v", err), w)
		return
	}

	// get the encoded password so not stored in plaintext in the database
	password, err := utils.SaltAndHashPassword(reqBody.Password)
	if err != nil {
		utils.LogAndAddServerError(err, w)
		return
	}

	// execute the query in the database
	wasErr := database.RegisterUser(reqBody.Username, reqBody.Email, password, w)
	if wasErr {
		return
	}

	// add new entry for preference tracker for the new user to manage next preference
	err = database.AddNewUserPreferenceTracker(reqBody.Username)
	if err != nil {
		utils.LogAndAddServerError(err, w)
		return
	}

	json.NewEncoder(w).Encode("user registered")
}
