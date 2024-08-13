package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lukelaurie/TikTokAutomation/backend/internal/database"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/middleware"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/model"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/utils"
)

func AddNewPreference(w http.ResponseWriter, r *http.Request) {
	// pull out the username from the cookie
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "Username not found in context", http.StatusInternalServerError)
		return
	}

	// get the preference passed into the request body 
	var reqBody model.Preference
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		panic(fmt.Errorf("request decode error: %v", err))
	}

	// first get the preference tracker from the database to determine what number to insert preference at
	preferenceTracker, err := database.RetrievePreferenceTracker(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// verify that user only has 10 preferences max
	if preferenceTracker.CurPreferenceCount >= 10 {
		http.Error(w, "user cannot have more than 10 preferences", http.StatusBadRequest)
		return
	} 
	
	// need to keep track of index in preference so can know the next preference to look at
	nextPreferenceIndex := preferenceTracker.CurPreferenceCount + 1 

	// try adding the new preference into the database
	err = database.AddNewUserPreference(reqBody, username, nextPreferenceIndex)
	if err != nil {
		utils.LogAndAddServerError(err, w)
		return
	}

	// increment the count in the database by one 
	err = database.IncrementPreferenceTracker(username, nextPreferenceIndex, true)
	if err != nil {
		utils.LogAndAddServerError(err, w)
		return
	}

	json.NewEncoder(w).Encode(username)
}