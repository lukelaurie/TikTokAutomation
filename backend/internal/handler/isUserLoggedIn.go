package handler

import (
	"encoding/json"
	"net/http"
)

func IsLoggedIn(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("valid")
}