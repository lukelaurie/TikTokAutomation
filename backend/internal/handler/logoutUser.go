package handler

import (
	"encoding/json"
	"net/http"
)

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// replace the old cookie with new one that about to expire
	http.SetCookie(w, &http.Cookie{
		Name:     "TikTokAutomationToken",
		Value:    "",
		MaxAge: -1,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	json.NewEncoder(w).Encode("user logged out")
}
