package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lukelaurie/TikTokAutomation/backend/internal/middleware"
)

func RetrievePreferences(w http.ResponseWriter, r *http.Request) {
	username, ok := middleware.GetUsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "Username not found in context", http.StatusInternalServerError)
        return
	}
	json.NewEncoder(w).Encode(username)
}