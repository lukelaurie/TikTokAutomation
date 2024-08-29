package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/handler"
	"github.com/lukelaurie/TikTokAutomation/backend/internal/middleware"
)

func InitializeRoutes(isTestMode bool) *mux.Router {
	router := mux.NewRouter()

	// create a subrouter for all routes to start with /api
	apiRouter := router.PathPrefix("/api").Subrouter()

	// public routes with no middleware
	apiRouter.HandleFunc("/register-user", handler.RegisterUser).Methods("POST")
	apiRouter.HandleFunc("/login", handler.LoginUser).Methods("GET")

	// private routes that require middleware
	protectedRouter := apiRouter.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(middleware.CheckAuthMiddleware) // apply the middleware to first authorize
	
	protectedRouter.HandleFunc("/upload-video", func(w http.ResponseWriter, r *http.Request) {
		handler.UploadVideo(isTestMode, w, r)
	}).Methods("POST")
	protectedRouter.HandleFunc("/retrieve-preferences", handler.RetrievePreferences).Methods("GET")
	protectedRouter.HandleFunc("/add-new-preference", handler.AddNewPreference).Methods("POST")

	return router
}
