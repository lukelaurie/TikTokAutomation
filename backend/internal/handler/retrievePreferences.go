package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lukelaurie/TikTokAutomation/backend/internal/database"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/middleware"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/utils"
)

func RetrievePreferences(w http.ResponseWriter, r *http.Request) {
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "Username not found in context", http.StatusInternalServerError)
		return
	}

	preferences, err := database.RetrieveAllUserPreferences(username)
	if err != nil {
		utils.LogAndAddServerError(err, w)
		return
	}

	json.NewEncoder(w).Encode(preferences)
}
