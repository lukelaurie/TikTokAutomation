package route

import (
	"github.com/gorilla/mux"
	uploadVideoHandler "github.com/lukelaurie/TikTokAutomation/backend/internal/handler"
)

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter() 
	router.HandleFunc("/api/upload-video", uploadVideoHandler.UploadVideo).Methods("POST")

	return router
}