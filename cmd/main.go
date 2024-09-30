package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/KarmaBeLike/SongLibrary/config"
	"github.com/KarmaBeLike/SongLibrary/internal/database"
	"github.com/KarmaBeLike/SongLibrary/internal/repository"
	"github.com/KarmaBeLike/SongLibrary/internal/routers"
	"github.com/KarmaBeLike/SongLibrary/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return
	}

	db, err := database.OpenDB(cfg)
	if err != nil {
		slog.Error("failed to connect to db", slog.Any("error", err))
		return
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		slog.Error("error running migrations", slog.Any("error", err))
		return
	}

	songRepo := repository.NewSongRepository(db)
	songService := service.NewSongService(songRepo, cfg.ExternalAPIURL)
	router := routers.SetupRoutes(songService)

	port := cfg.Port
	log.Printf("Server is running on port %d...", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
