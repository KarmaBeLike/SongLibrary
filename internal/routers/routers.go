package routers

import (
	"github.com/gorilla/mux"

	"github.com/KarmaBeLike/SongLibrary/internal/handlers"
)

// SetupRoutes создает маршруты для приложения с использованием Gorilla Mux
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/songs", handlers.SongsHandler).Methods("GET")
	router.HandleFunc("/api/songs/{id}", handlers.GetSongHandler).Methods("GET")
	router.HandleFunc("/api/songs/{id}", handlers.DeleteSongHandler).Methods("DELETE")
	router.HandleFunc("/api/songs/{id}", handlers.UpdateSongHandler).Methods("PUT")
	router.HandleFunc("/api/songs", handlers.AddSongHandler).Methods("POST")

	return router
}
