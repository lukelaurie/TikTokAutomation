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
	var reqBody model.RegisterRequest

	// Decode the body of the request into the struct
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		panic(fmt.Errorf("request decode error: %v", err))
	}

	// get the encoded passworded so not stored in plaintext in the database
	password, err := utils.SaltAndHashPassword(reqBody.Password)
	if err != nil {
		panic(err)
	}

	// execute the query in the database
	wasErr := database.RegisterUser(reqBody.Username, reqBody.Email, password, w)
	if wasErr {
		return
	}

	json.NewEncoder(w).Encode("user regiserd")
}
