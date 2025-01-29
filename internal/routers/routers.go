package routers

import (
	"github.com/KarmaBeLike/SongLibrary/internal/handlers"
	"github.com/KarmaBeLike/SongLibrary/internal/service"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(songService *service.SongService) *mux.Router {
	songHandler := handlers.NewSongClient(songService)

	router := mux.NewRouter()
	router.HandleFunc("/api/songs", songHandler.GetSongs).Methods("GET")
	router.HandleFunc("/api/songs/lyrics", songHandler.GetSongLyrics).Methods("GET")
	router.HandleFunc("/api/songs", songHandler.DeleteSong).Methods("DELETE")
	router.HandleFunc("/api/songs", songHandler.UpdateSong).Methods("PATCH")
	router.HandleFunc("/api/songs", songHandler.AddSong).Methods("POST")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}
