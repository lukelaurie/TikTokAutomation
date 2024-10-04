package middleware

import (
	"net/http"
	"os"
)

func CheckAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		correctApiKey := os.Getenv("xApiKey")

		if apiKey != correctApiKey {
			http.Error(w, "Forbidden: Invalid API key", http.StatusForbidden)
			return
		}
	
		next.ServeHTTP(w, r)
	})
}