package route

import (
	"github.com/gorilla/mux"
	handler "github.com/lukelaurie/TikTokAutomation/backend/internal/handler"
)

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter() 
	router.HandleFunc("/api/upload-video", handler.UploadVideo).Methods("POST")
	router.HandleFunc("/api/register-user", handler.RegisterUser).Methods("POST")
	router.HandleFunc("/api/login", handler.LoginUser).Methods("POST")

	return router
}